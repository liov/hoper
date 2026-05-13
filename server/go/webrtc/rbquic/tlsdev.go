package rbquic

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"math/big"

	"github.com/quic-go/quic-go"
)

// NextProto 与 DevClientTLS / DevServerTLS 一致，用于远程浏览数据面调试。
const NextProto = "hoper-rb"

// DevServerTLS 生成临时自签证书，仅用于内网/开发；生产须使用正式证书。
func DevServerTLS() *tls.Config {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{{Certificate: [][]byte{certDER}, PrivateKey: priv}},
		NextProtos:   []string{NextProto},
	}
}

// DevClientTLS 对应 DevServerTLS，开启跳过校验（仅调试）。
func DevClientTLS() *tls.Config {
	return &tls.Config{InsecureSkipVerify: true, NextProtos: []string{NextProto}}
}

// ListenDev 在 addr 上监听 QUIC（如 "127.0.0.1:0"）。
func ListenDev(addr string) (*quic.EarlyListener, error) {
	return quic.ListenAddrEarly(addr, DevServerTLS(), &quic.Config{})
}

// DialDev 连接对端 QUIC（IPv6 直连路径可传 `[ipv6]:port`）。
func DialDev(ctx context.Context, addr string) (*quic.Conn, error) {
	return quic.DialAddrEarly(ctx, addr, DevClientTLS(), &quic.Config{})
}
