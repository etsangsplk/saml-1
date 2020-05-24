package saml_test

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ucarion/dsig"
	"github.com/ucarion/saml"
)

// The cert and key used in these tests were generated by running:
//
// openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365
// -nodes -subj "/C=US/ST=Oregon/L=Portland/O=Company
// Name/OU=Org/CN=www.example.com"
//
// Which generated the certificate used in TestVerify. It also generated this
// corresponding RSA private key:
//
// -----BEGIN PRIVATE KEY-----
// MIIJQgIBADANBgkqhkiG9w0BAQEFAASCCSwwggkoAgEAAoICAQC4yrnAdbSwMQz8
// CPL0ir2TbhNadYCjukgxQxxjGeVVuRxHFzKyLlaM17ZG6x67F7DqvxXwO8Bvfnur
// eFDG8uq13eIni3Tvl38G4pFQOLHgY/hvKS9E4L8bBPDaXag6jD463arduvOo0ZzB
// movgoqw4ryze53I2z8D+eIPa3WmRX93gZ69UemcAo7MJv/WzpByimefah/+TSTFd
// r5q1kVOxaJclcoAbTcDRxZEVGc6l/BRZQHrz5VG3JEO8M1KuzAq0nMyFkyP1ynGv
// P66syV+DYcqjSJ3cB+fgeIC1VT2ZgwzdxBWJeRIgCNT8YOmI9FpiYrurEcU9iLnf
// HQ0O3422QLZrJVxyY47c7q9kBPUNUZ+ewWgPeSeRDSl3L0wfknNebemmZr7IRYUV
// h6bTIqB/MRORp70E/Vnp5QBrofdkaF8Z6zslnz50FAj3pZVsmEaQjh/yBqxghV4H
// qivNmf+6p8/38PjhlU9x/M1dmmjr1NKau4ILrS9meEL2+rweqV4TZq++T7d57Bor
// oz1zytKXp/QwvGQcpD+TRtuhHZkpLHROyGSdCR9XnYBI5FeQOVi53vhT+BDrbrXE
// +CkfmALylLFsflBeF/9IhmyHAcEpkj/0FZYjAXoEZEEeAphawd/qVspt30J7X35O
// cqTIPDilNVL+WWMi0lHsWO5h9FzPMwIDAQABAoICAH3d353ezp8AGgcFlW7JnYzw
// +g+wX1mmBYxAWPKLbfDwr/kgLPC+rUcrmsU9WuY2odOTKj9Cg7WtolDOF78bMJGF
// u4gR7ilPuD8ZTb8ljsr3bP1SQRcaOjEOMXubNX4DjlOMLtjugQ6pD6uzN7lfNA08
// DEUbwmjhI2Rw8+a8zy4s7TTvirXw1X3TAp0Oei3NB5AdYpYv8f4BabWVabxoa2g4
// hFMGZYmzcTWw6zxDIsVeKQIN8HF17i3fbp+fGZ9j7ZrN/mSxL1o4dSzYJIMeeodD
// scF8McHwRJlZmtloYRfR8o6PA9hqddUKDwCEhi05uuKuu4MvDHj4SxpUcFOEI8Iq
// cxtM4m8kbLPpWnHYC0xUUigP5lY0P8HoQqO6ktLJNv2s2pEKxXFCA7kTiA2xyiEC
// iisuoyGGlXQ8HB8L1ShYDrwjPlbJ7CCTl3yZ8kruq6kp5kF/Cps39cmM0WsvJXK0
// OjYBNFAO2RnoRaoJQrh60rWPjK/JoqpUGrRPG7+k1VuGdGBaOw4YvC339koALXpK
// sEDvwIztPsv2AwE0WMw5NmhnRtmNUXb3LhlrjPOZ5e2HyRhdrLpBu5y79tcF8VN9
// mksfIQqhNxKpidBkju4u6nyYOvQL013vKJSJRfXxYXRKQgGO0gBhasjjwBRjeOvY
// uLKVpnJ1Ncq4nr1yFLdZAoIBAQDfl0kxOJeFBYSVQRKxiXIL+xXm9mSQ3fu2Ir4T
// JoMkkz0pBkXCI5aC+JnR8CmmTVXD+T5BDYcXGngGAZQZDPvFhXs9J0dpJJRVajBT
// Cj7OPRu51o95gjmYVpOHtKwRQal8LerPajd60YsRNNuMsExkJkmirUly3bkaevEh
// Gqt0RK6/qTGDy5M0u5KYhMy/mgg+mSWLfjLmdR7FaszVmIEa6WZr5dTMmmidPq7b
// GHWvXwLZg6VIeGzRhAb33XLDBrB/S/IOY98dz4hasVZ/f+EskgFqkEXe/X1i1YSN
// PT1VNSOWA6Rj/h+rmA7NOdLZMrzi12sId2Y6sjvbPERNPhftAoIBAQDTk7xn6usz
// /MIJk4HhYwlVJadPnfsFP7Z0AGXGY83OclpsiDpJPC8nWGvc40oqtmqz4cOcdy+H
// KFtu1JGyDn9W36y9+NViQp+RJ1NKosdV4/N7L9nXXi1y7uNe9QvdTjtS44xXFEk5
// FLDoKzGDXkPp77eA6BPfqMMFymf8mgq+MpWioKLiR43w+Zc+/Ncz6zMSsr/nPH6e
// 1Gjh0Nva4/M5aelJU+i5P1bJlcrRs6//N3RQjPgCBF5NDj2SEseAH8cQkxhfsXB+
// xfWyY7ocGPNlO+sGarLqaftqSSD1J7wZ8dbgHysnTJkdhhWJmRbYxSoZBwSG1CSI
// kvDugRZ8N1+fAoIBACCb2crZ7A80bM+vu+A0oXNp3RngGW6fUVSQ4JO+bCXra2IO
// TiIwOoVDaHubwRdF9BouwYuPQ4J1E8gcdtLod9eozf5vOhT1hsSmRgH2Xo6Jjv+d
// cTNRcMDs73s9OFMT9nnr4HD7lrfM07FguhxcoeeBRf/5sdqUx6g7AevIDfVZBvtg
// 253TFNb9/DVOOOZAuq8WeslLUHUX47L7DoCgS0P3gj5+OHjWlCdKuwmtGYzIGIxM
// jNBy77vmu3Vu0Ivs79TA6L58hk+8srA3aNwTdG2hpZ87B1WsNpsxdLF8mvNQWq5I
// PbNvnoLSHGaF5mBS7AVRUYTclQY+dEhXE8cIJUkCggEBAMU9bN7TugD1GU8kHGip
// kwG14IvwkxsJkmYCGN8iG7LiGDolpXCwkqTzYVrC6Vl4RXD8fwdWdRBjJxnjQQ/l
// RAEQ9FEFsKexxF/lcVia94mywEGPEl4chfInkf/sIetmCxfy2do0Jy73gxRtb/Mv
// 5dAokcGymRRgl67GSrrKQEmfjq/VYQPiAQktJTqrK1RTZ4F+8jf3xXL8QeqCcvNU
// nmJfwgOCHerUiWvUIQfto50hbWXKhUocGG1tYSjUKPfgqAtjlc1f9ae5lJuBLPcU
// q5MskKWiwriVpLQpCHiDWnA1bEPzyp8QYY2MeneUKCBdbil2yVmIW6aWldVCsluK
// o7ECggEAbp+MZOPzKYTEGVWLNQh0CVairBVrOexlOFrup7sOW0NFQXu8ExHsRsgC
// HMEvBj24jJM6FeaJ4Fkc1WAfJqY0KnpWeEPFzLY9W7ZEHbkyiHJ0DzvReXPYGWSC
// Qj0dgv0jfDODdsfTqI6zW/WXHEQ8399JiAEVGVphMUo2oY+rhDAiZCFFlt7heyq2
// fLf4MAmc3vK6slbyaDb9kYm+fsiCBVqvwVIKvIZ1/IOOU5q6KQIYjJXryLIBORuw
// 3jlAmnFMZFC0dBPJAHeon8m47S/1Te2EkyH1D1GvcDnE07PjhFUl3LpbD4qrw0Wv
// tRNOxnQnlHJKcCgbfcUOD3hpFKtY9g==
// -----END PRIVATE KEY-----

func TestVerify(t *testing.T) {
	block, _ := pem.Decode([]byte(`-----BEGIN CERTIFICATE-----
MIIFXDCCA0QCCQCl4WZtbTlavDANBgkqhkiG9w0BAQsFADBwMQswCQYDVQQGEwJV
UzEPMA0GA1UECAwGT3JlZ29uMREwDwYDVQQHDAhQb3J0bGFuZDEVMBMGA1UECgwM
Q29tcGFueSBOYW1lMQwwCgYDVQQLDANPcmcxGDAWBgNVBAMMD3d3dy5leGFtcGxl
LmNvbTAeFw0yMDA1MjAxNzI0MzFaFw0yMTA1MjAxNzI0MzFaMHAxCzAJBgNVBAYT
AlVTMQ8wDQYDVQQIDAZPcmVnb24xETAPBgNVBAcMCFBvcnRsYW5kMRUwEwYDVQQK
DAxDb21wYW55IE5hbWUxDDAKBgNVBAsMA09yZzEYMBYGA1UEAwwPd3d3LmV4YW1w
bGUuY29tMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAuMq5wHW0sDEM
/Ajy9Iq9k24TWnWAo7pIMUMcYxnlVbkcRxcysi5WjNe2Ruseuxew6r8V8DvAb357
q3hQxvLqtd3iJ4t075d/BuKRUDix4GP4bykvROC/GwTw2l2oOow+Ot2q3brzqNGc
wZqL4KKsOK8s3udyNs/A/niD2t1pkV/d4GevVHpnAKOzCb/1s6Qcopnn2of/k0kx
Xa+atZFTsWiXJXKAG03A0cWRFRnOpfwUWUB68+VRtyRDvDNSrswKtJzMhZMj9cpx
rz+urMlfg2HKo0id3Afn4HiAtVU9mYMM3cQViXkSIAjU/GDpiPRaYmK7qxHFPYi5
3x0NDt+NtkC2ayVccmOO3O6vZAT1DVGfnsFoD3knkQ0pdy9MH5JzXm3ppma+yEWF
FYem0yKgfzETkae9BP1Z6eUAa6H3ZGhfGes7JZ8+dBQI96WVbJhGkI4f8gasYIVe
B6orzZn/uqfP9/D44ZVPcfzNXZpo69TSmruCC60vZnhC9vq8HqleE2avvk+3eewa
K6M9c8rSl6f0MLxkHKQ/k0bboR2ZKSx0TshknQkfV52ASORXkDlYud74U/gQ6261
xPgpH5gC8pSxbH5QXhf/SIZshwHBKZI/9BWWIwF6BGRBHgKYWsHf6lbKbd9Ce19+
TnKkyDw4pTVS/lljItJR7FjuYfRczzMCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEA
r6UAa9n4FkiA4ZqugCJEoC5Ehc1X/qdNFkY4EIHc33sqscqVZhHC0MbfNmKuiirk
XKTR+M3U62IvD8HXpkBMTYMpnvsH4jFuP3SpTFfUuqarueqsawiPAejhjF9829fg
K1+s1rD/fI3H3UuHWChTXKA4KpnCYr5B1om4ZoCcTVVdZjhO256iM7p/DHze08Eo
Rdhaj+rgs6NC5vLHWX9bezACeqA3YwJYHRH0zuoCQfRKXkikIjj18wpWNARFhDoQ
FEhJXIAO/skpuK6Q9Ml1wWuFaqgXtKN1iVzuGi7P8O3bCLexwmqnmsnEZPPpzjoQ
T8zVIjCH6jBX533f1B745IrGNzMSr6YC/9RT3DrPoNT9pCAozSoZxldqIegxLgWG
zBT6jj/fR92E5kJh8Hy3koeXGkyAkcHB0PH8yyFtYIlP0stENkG/fDCLuMUqf6GZ
P/oSyJH1Ro/qV6kwc1XYDB+6NGC8Xd1JQKZD49c/GZYpo77ZYKQtCoTrMuPKSG5/
jP7OTrdylTj+V4r7jYLLpvWCUe0ON0QPKClo+15tXATWep6PFk0U5W+efvavG70e
Fu9GKMOkTgv5F/ngzDgXKo7T6poRDZAgolUAq2kwDUp42AVx/7UqmOdp0yUTNmJG
A70UwPLAvWk5vX1IMpaEFjBd3LqWLeSmbKZ03zr1jnA=
-----END CERTIFICATE-----`))

	cert, err := x509.ParseCertificate(block.Bytes)
	assert.NoError(t, err)

	now, err := time.Parse(time.RFC3339, "2020-05-23T01:46:00Z")
	assert.NoError(t, err)

	type testCase struct {
		Name     string
		Response saml.Response
		Error    error
	}

	testCases := []testCase{
		testCase{
			Name:  "unsigned",
			Error: saml.ErrResponseNotSigned,
		},
		testCase{
			Name:  "invalid_signature",
			Error: rsa.ErrVerification,
		},
		testCase{
			Name:  "wrong_issuer",
			Error: saml.ErrInvalidIssuer,
		},
		testCase{
			Name:  "wrong_recipient",
			Error: saml.ErrInvalidRecipient,
		},
		testCase{
			Name:  "before_conditions_not_before",
			Error: saml.ErrAssertionExpired,
		},
		testCase{
			Name:  "after_conditions_not_on_or_after",
			Error: saml.ErrAssertionExpired,
		},
		testCase{
			Name:  "after_subject_confirmation_data_not_on_or_after",
			Error: saml.ErrAssertionExpired,
		},
		testCase{
			Name: "valid",
			Response: saml.Response{
				XMLName: xml.Name{Space: "urn:oasis:names:tc:SAML:2.0:protocol", Local: "Response"},
				Assertion: saml.Assertion{
					XMLName: xml.Name{Space: "urn:oasis:names:tc:SAML:2.0:assertion", Local: "Assertion"},
					Issuer: saml.Issuer{
						XMLName: xml.Name{Space: "urn:oasis:names:tc:SAML:2.0:assertion", Local: "Issuer"},
						Name:    "alice",
					},
					Subject: saml.Subject{
						XMLName: xml.Name{Space: "urn:oasis:names:tc:SAML:2.0:assertion", Local: "Subject"},
						NameID: saml.NameID{
							XMLName: xml.Name{Space: "urn:oasis:names:tc:SAML:2.0:assertion", Local: "NameID"},
							Format:  "format",
							Value:   "jdoe@example.com",
						},
						SubjectConfirmation: saml.SubjectConfirmation{
							XMLName: xml.Name{Space: "urn:oasis:names:tc:SAML:2.0:assertion", Local: "SubjectConfirmation"},
							Method:  "method",
							SubjectConfirmationData: saml.SubjectConfirmationData{
								XMLName:      xml.Name{Space: "urn:oasis:names:tc:SAML:2.0:assertion", Local: "SubjectConfirmationData"},
								NotOnOrAfter: now.Add(time.Minute),
								Recipient:    "bob",
							},
						},
					},
					Conditions: saml.Conditions{
						XMLName:      xml.Name{Space: "urn:oasis:names:tc:SAML:2.0:assertion", Local: "Conditions"},
						NotBefore:    now.Add(-time.Minute),
						NotOnOrAfter: now.Add(time.Minute),
					},
					AttributeStatement: saml.AttributeStatement{
						XMLName: xml.Name{Space: "urn:oasis:names:tc:SAML:2.0:assertion", Local: "AttributeStatement"},
						Attributes: []saml.Attribute{
							saml.Attribute{
								XMLName:    xml.Name{Space: "urn:oasis:names:tc:SAML:2.0:assertion", Local: "Attribute"},
								Name:       "attr1",
								NameFormat: "fmt1",
								Value:      "value1",
							},
							saml.Attribute{
								XMLName:    xml.Name{Space: "urn:oasis:names:tc:SAML:2.0:assertion", Local: "Attribute"},
								Name:       "attr2",
								NameFormat: "fmt2",
								Value:      "value2",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s.xml", tt.Name))
			assert.NoError(t, err)

			res, err := saml.Verify(base64.StdEncoding.EncodeToString(b), "alice", cert, "bob", now)
			res.Signature = dsig.Signature{} // clear this out to make test cases a bit less verbose

			assert.Equal(t, tt.Response, res)
			assert.Equal(t, tt.Error, err)
		})
	}
}

func TestGetEntityIDCertificateAndRedirectURL(t *testing.T) {
	block, _ := pem.Decode([]byte(`-----BEGIN CERTIFICATE-----
MIIFXDCCA0QCCQCl4WZtbTlavDANBgkqhkiG9w0BAQsFADBwMQswCQYDVQQGEwJV
UzEPMA0GA1UECAwGT3JlZ29uMREwDwYDVQQHDAhQb3J0bGFuZDEVMBMGA1UECgwM
Q29tcGFueSBOYW1lMQwwCgYDVQQLDANPcmcxGDAWBgNVBAMMD3d3dy5leGFtcGxl
LmNvbTAeFw0yMDA1MjAxNzI0MzFaFw0yMTA1MjAxNzI0MzFaMHAxCzAJBgNVBAYT
AlVTMQ8wDQYDVQQIDAZPcmVnb24xETAPBgNVBAcMCFBvcnRsYW5kMRUwEwYDVQQK
DAxDb21wYW55IE5hbWUxDDAKBgNVBAsMA09yZzEYMBYGA1UEAwwPd3d3LmV4YW1w
bGUuY29tMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAuMq5wHW0sDEM
/Ajy9Iq9k24TWnWAo7pIMUMcYxnlVbkcRxcysi5WjNe2Ruseuxew6r8V8DvAb357
q3hQxvLqtd3iJ4t075d/BuKRUDix4GP4bykvROC/GwTw2l2oOow+Ot2q3brzqNGc
wZqL4KKsOK8s3udyNs/A/niD2t1pkV/d4GevVHpnAKOzCb/1s6Qcopnn2of/k0kx
Xa+atZFTsWiXJXKAG03A0cWRFRnOpfwUWUB68+VRtyRDvDNSrswKtJzMhZMj9cpx
rz+urMlfg2HKo0id3Afn4HiAtVU9mYMM3cQViXkSIAjU/GDpiPRaYmK7qxHFPYi5
3x0NDt+NtkC2ayVccmOO3O6vZAT1DVGfnsFoD3knkQ0pdy9MH5JzXm3ppma+yEWF
FYem0yKgfzETkae9BP1Z6eUAa6H3ZGhfGes7JZ8+dBQI96WVbJhGkI4f8gasYIVe
B6orzZn/uqfP9/D44ZVPcfzNXZpo69TSmruCC60vZnhC9vq8HqleE2avvk+3eewa
K6M9c8rSl6f0MLxkHKQ/k0bboR2ZKSx0TshknQkfV52ASORXkDlYud74U/gQ6261
xPgpH5gC8pSxbH5QXhf/SIZshwHBKZI/9BWWIwF6BGRBHgKYWsHf6lbKbd9Ce19+
TnKkyDw4pTVS/lljItJR7FjuYfRczzMCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEA
r6UAa9n4FkiA4ZqugCJEoC5Ehc1X/qdNFkY4EIHc33sqscqVZhHC0MbfNmKuiirk
XKTR+M3U62IvD8HXpkBMTYMpnvsH4jFuP3SpTFfUuqarueqsawiPAejhjF9829fg
K1+s1rD/fI3H3UuHWChTXKA4KpnCYr5B1om4ZoCcTVVdZjhO256iM7p/DHze08Eo
Rdhaj+rgs6NC5vLHWX9bezACeqA3YwJYHRH0zuoCQfRKXkikIjj18wpWNARFhDoQ
FEhJXIAO/skpuK6Q9Ml1wWuFaqgXtKN1iVzuGi7P8O3bCLexwmqnmsnEZPPpzjoQ
T8zVIjCH6jBX533f1B745IrGNzMSr6YC/9RT3DrPoNT9pCAozSoZxldqIegxLgWG
zBT6jj/fR92E5kJh8Hy3koeXGkyAkcHB0PH8yyFtYIlP0stENkG/fDCLuMUqf6GZ
P/oSyJH1Ro/qV6kwc1XYDB+6NGC8Xd1JQKZD49c/GZYpo77ZYKQtCoTrMuPKSG5/
jP7OTrdylTj+V4r7jYLLpvWCUe0ON0QPKClo+15tXATWep6PFk0U5W+efvavG70e
Fu9GKMOkTgv5F/ngzDgXKo7T6poRDZAgolUAq2kwDUp42AVx/7UqmOdp0yUTNmJG
A70UwPLAvWk5vX1IMpaEFjBd3LqWLeSmbKZ03zr1jnA=
-----END CERTIFICATE-----`))

	expectedCert, err := x509.ParseCertificate(block.Bytes)
	assert.NoError(t, err)

	b, err := ioutil.ReadFile("tests/valid_idp_metadata.xml")
	assert.NoError(t, err)

	var metadata saml.EntityDescriptor
	assert.NoError(t, xml.Unmarshal(b, &metadata))

	fmt.Println(metadata)

	entityID, cert, location, err := metadata.GetEntityIDCertificateAndRedirectURL()
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", entityID)
	assert.True(t, expectedCert.Equal(cert))
	assert.Equal(t, "https://example.com/saml/redirect", location.String())
}