// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package node

import (
	"autohost-cli/internal/adapters/agentconfig"
	"autohost-cli/internal/adapters/enrollapi"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

// EnrollService handles node enrollment with the AutoHost cloud API.
type EnrollService struct{}

func (s *EnrollService) Link(api, token, name string) error {
	nd := s.GatherNodeData()
	if name != "" {
		nd.HostName = name
	}

	csrPEM, keyPEM, err := generateNodeCSR(nd.HostName)
	if err != nil {
		return fmt.Errorf("error generando CSR para mTLS: %w", err)
	}

	req := NodeRequest{
		EnrollToken:  strings.TrimSpace(token),
		HostName:     strings.TrimSpace(nd.HostName),
		IPLocal:      strings.TrimSpace(nd.IPLocal),
		OS:           strings.TrimSpace(nd.OS),
		Arch:         strings.TrimSpace(nd.Arch),
		VersionAgent: strings.TrimSpace(nd.VersionAgent),
		CSRPEM:       csrPEM,
	}

	client := enrollapi.NewUserClient(api, token)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var resp NodeResponse
	_, err = client.PostJSON(ctx, "/v1/enrollments/enroll", req, &resp)
	if err != nil {
		return fmt.Errorf("enrollment fallido: %w", err)
	}
	fmt.Println("💾 Guardando configuración y certificados mTLS...")

	clientKey := keyPEM
	if resp.ClientKey != "" {
		clientKey = resp.ClientKey
	}

	if err := agentconfig.Save(agentconfig.AgentConfig{
		ApiToken:     resp.ApiToken,
		RefreshToken: resp.RefreshToken,
		ApiURL:       api,
		NodeID:       resp.NodeID,
		GRPCAddress:  resp.GRPCAddress,
		CACert:       resp.CACert,
		ClientCert:   resp.ClientCert,
		ClientKey:    clientKey,
	}); err != nil {
		return fmt.Errorf("error guardando configuración: %w", err)
	}

	fmt.Println("✅ Configuración actualizada en /etc/autohost/config.yaml")
	fmt.Println()
	fmt.Println("📝 Próximos pasos:")
	fmt.Println("  1. Habilitar servicio: sudo systemctl enable autohost-agent")
	fmt.Println("  2. Iniciar servicio:   sudo systemctl start autohost-agent")
	fmt.Println("  3. Verificar estado:   sudo systemctl status autohost-agent")
	return nil
}

func (s *EnrollService) GatherNodeData() *NodeData {
	nd := &NodeData{}
	if out, err := exec.Command("hostname").Output(); err == nil {
		nd.HostName = strings.TrimSpace(string(out))
	}
	if out, err := exec.Command("uname", "-o").Output(); err == nil {
		nd.OS = strings.TrimSpace(string(out))
	}
	if out, err := exec.Command("uname", "-m").Output(); err == nil {
		nd.Arch = strings.TrimSpace(string(out))
	}
	nd.IPLocal = DetectLocalIP()
	nd.VersionAgent = "0.1.0"
	return nd
}

func DetectLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

func generateNodeCSR(hostname string) (csrPEM string, keyPEM string, err error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	cn := strings.TrimSpace(hostname)
	if cn == "" {
		cn = "autohost-node"
	}

	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: []string{"Autohost Node"},
		},
	}

	csrDER, err := x509.CreateCertificateRequest(rand.Reader, template, priv)
	if err != nil {
		return "", "", err
	}

	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return "", "", err
	}

	csrBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})
	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	return string(csrBytes), string(keyBytes), nil
}
