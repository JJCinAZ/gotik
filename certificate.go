package gotik

import (
	"fmt"
	"time"
)

type Certificate struct {
	ID              string        `json:".id"`
	Trusted         bool          `json:"trusted"`
	Expired         bool          `json:"expired"`
	Revoked         bool          `json:"revoked"`
	Issued          bool          `json:"issued"`
	Authority       bool          `json:"authority"`
	CRL             bool          `json:"crl"`
	SmartCardKey    bool          `json:"smart-card-key"`
	PrivateKey      bool          `json:"private-key"`
	Name            string        `json:"name"`
	Issuer          string        `json:"issuer"`
	DigestAlgorithm string        `json:"digest-algorithm"`
	KeyType         string        `json:"key-type"`
	KeySize         int           `json:"key-size"`
	Country         string        `json:"country"`
	Organization    string        `json:"organization"`
	CN              string        `json:"common-name"`
	SAN             string        `json:"subject-alt-name"`
	DaysValid       int           `json:"days-valid"`
	KeyUsage        string        `json:"key-usage"`
	Serial          string        `json:"serial-number"`
	Fingerprint     string        `json:"fingerprint"`
	InvalidBefore   time.Time     `json:"invalid-before"`
	InvalidAfter    time.Time     `json:"invalid-after"`
	ExpiresAfter    time.Duration `json:"expires-after"`
	// akid=f0126d62cbd3bd39811b682d58fc591a9d35ea22
	// skid=5108b87b0e51b690fe253b3890b270c00089c2cd
}

type CertImportResults struct {
	CertificatesImported   int `json:"certificates-imported"`
	PrivateKeysImported    int `json:"private-keys-imported"`
	FilesImported          int `json:"files-imported"`
	DecryptionFailures     int `json:"decryption-failures"`
	KeysWithNoCertificates int `json:"keys-with-no-certificates"`
}

func parseCertificate(props map[string]string) Certificate {
	entry := Certificate{
		ID:              props[".id"],
		Name:            props["name"],
		Issuer:          props["issuer"],
		DigestAlgorithm: props["digest-algorithm"],
		KeyType:         props["key-type"],
		KeySize:         parseInt(props["key-size"]),
		Country:         props["country"],
		Organization:    props["organization"],
		CN:              props["common-name"],
		SAN:             props["subject-alt-name"],
		DaysValid:       parseInt(props["days-valid"]),
		Trusted:         parseBool(props["trusted"]),
		Expired:         parseBool(props["expired"]),
		Revoked:         parseBool(props["revoked"]),
		Issued:          parseBool(props["issued"]),
		Authority:       parseBool(props["authority"]),
		CRL:             parseBool(props["crl"]),
		SmartCardKey:    parseBool(props["smart-card-key"]),
		PrivateKey:      parseBool(props["private-key"]),
		KeyUsage:        props["key-usage"],
		Serial:          props["serial-number"],
		Fingerprint:     props["fingerprint"],
		InvalidBefore:   parseTime(props["invalid-before"]),
		InvalidAfter:    parseTime(props["invalid-after"]),
		ExpiresAfter:    parseDuration(props["expires-after"]),
	}
	return entry
}

func parseCertImportResults(props map[string]string) CertImportResults {
	entry := CertImportResults{
		CertificatesImported:   parseInt(props["certificates-imported"]),
		PrivateKeysImported:    parseInt(props["private-keys-imported"]),
		FilesImported:          parseInt(props["files-imported"]),
		DecryptionFailures:     parseInt(props["decryption-failures"]),
		KeysWithNoCertificates: parseInt(props["keys-with-no-certificates"]),
	}
	return entry
}

func (c *Certificate) String() string {
	return fmt.Sprintf("%s: cn=%s, issuer=%s, key=%s(%d), san=%s, valid: %d, trusted/expired/revoked=%t/%t/%t",
		c.Name, c.CN, c.Issuer, c.KeyType, c.KeySize, c.SAN, c.DaysValid, c.Trusted, c.Expired, c.Revoked)
}

func (c *Client) CertificateImport(name, filename, passphrase string) (CertImportResults, error) {
	var r CertImportResults
	parts := make([]string, 0)
	parts = append(parts, "/certificate/import")
	parts = append(parts, fmt.Sprintf("=file-name=%s", filename))
	parts = append(parts, fmt.Sprintf("=passphrase=%s", passphrase))
	if c.majorVersion > 6 || (c.majorVersion == 6 && c.minorVersion >= 47) {
		if len(name) > 0 {
			parts = append(parts, fmt.Sprintf("=name=%s", name))
		}
	}
	detail, err := c.Run(parts...)
	if err == nil {
		r = parseCertImportResults(detail.Re[0].Map)
	}
	return r, err
}

func (c *Client) certPrint(parms ...string) ([]Certificate, error) {
	entries := make([]Certificate, 0)
	detail, err := c.RunCmd("/certificate/print", parms...)
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parseCertificate(detail.Re[i].Map))
		}
	}
	return entries, nil
}

func (c *Client) SetCertificateName(id string, name string) error {
	_, err := c.Run("/certificate/set", "=.id="+id, "=name="+name)
	return err
}

func (c *Client) GetCertificates() ([]Certificate, error) {
	return c.certPrint()
}

func (c *Client) RemoveCertificate(id string) error {
	_, err := c.Run("/certificate/remove", "=.id="+id)
	return err
}
