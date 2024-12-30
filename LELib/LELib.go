// library to create and check LE Account
// author: prr, azul software
// date: 29 Dec 2024
// copyright 2024 prr, azulsoftware
//
// 

package LELib

import (
	"fmt"
	"os"
    "time"
    "context"
//    "strings"
//    "net"
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509"
//    "crypto/x509/pkix"
//    "encoding/asn1"
    "encoding/pem"

   "golang.org/x/crypto/acme"
    yaml "github.com/goccy/go-yaml"
)

const LEProdUrl = "https://acme-v02.api.letsencrypt.org/directory"
const LETestUrl = "https://acme-staging-v02.api.letsencrypt.org/directory"

type AccountInfo struct {
    Contacts []string `yaml:"Contacts"`
}

type LEAcnt struct {
    AcntNam string `yaml:"AcntName"`
    AcntId string `yaml:"AcntId"`
    PrivKeyFilnam string `yaml:"PrivKeyFilnam"`
    PubKeyFilnam string `yaml:"PubKeyFilnam"`
    Updated time.Time `yaml:"update"`
    Contacts []string `yaml:"contacts"`
	Prod bool `yaml:"Prod"`
    LEUrl string `yaml:"LEUrl"`
}

type LELib struct {
	Name string
	LEDir string
	Dbg bool
    Prod bool
	Ctx context.Context
//    Client *acme.Client
//    LEAccount *acme.Account
}

func InitLELib(acntNam string, typ bool) (*LELib, error){

	var LELibObj  LELib

	if len(acntNam) == 0 {return nil, fmt.Errorf("no account name!")}
	LELibObj.Name = acntNam

	LEDir := os.Getenv("LEDir")
	if len(LEDir) == 0 {return nil, fmt.Errorf("cannot find LEDir!")}
	LELibObj.LEDir = LEDir
	LELibObj.Prod = typ

	return &LELibObj, nil
}

func (LELibObj *LELib) CreateLEAccount() (err error) {

    dbg := LELibObj.Dbg
	ctx := LELibObj.Ctx

	LEDir := os.Getenv("LEDir")
	if len(LEDir) == 0 {return fmt.Errorf("cannot find LEDir!")}

// contruct LE account name
	acntFilnam :=""
	privKeyFilnam := ""
	pubKeyFilnam := ""

	if LELibObj.Prod {
    	acntFilnam =  LEDir + "/" + LELibObj.Name + "LEProd.yaml"
	    privKeyFilnam = LEDir + "/" + LELibObj.Name + "LEProdPriv.key"
    	pubKeyFilnam = LEDir + "/" + LELibObj.Name + "LEProdPub.key"
	} else {
    	acntFilnam =  LEDir + "/" + LELibObj.Name + "LETest.yaml"
	    privKeyFilnam = LEDir + "/" + LELibObj.Name + "LETestPriv.key"
    	pubKeyFilnam = LEDir + "/" + LELibObj.Name + "LETestPub.key"
	}
    acntInfoFilnam := LEDir + "/" + LELibObj.Name + "_info.yaml"

    if dbg {
		fmt.Println("*************** dbg info ***************")
		fmt.Printf("account file:    %s\n", acntFilnam)
		fmt.Printf("account info:    %s\n", acntInfoFilnam)
        fmt.Printf("priv Key Filnam: %s\n", privKeyFilnam)
        fmt.Printf("pub Key Filnam:  %s\n", pubKeyFilnam)
		fmt.Println("*********** end dbg info ***************")
    }

    // check for existing key and yaml files
	if _, err := os.Stat(acntFilnam); err == nil {
		return fmt.Errorf("account file already exists!")
	}
	if _, err := os.Stat(privKeyFilnam); err == nil {
		os.Remove(privKeyFilnam)
		fmt.Printf("no account file -- removing private key file!\n")
	}
	if _, err := os.Stat(pubKeyFilnam); err == nil {
		os.Remove(privKeyFilnam)
		fmt.Printf("no account file -- removing public key file!\n")
	}

    acntData, err := os.ReadFile(acntInfoFilnam)
    if err != nil {return fmt.Errorf("account info: %v", err)}

    acntInfo := AccountInfo{}
    err = yaml.Unmarshal(acntData, &acntInfo)
    if err != nil {return fmt.Errorf("yaml Unmarshal account info: %v\n", err)}

//	if actInfo.Name != LELibObj.Name {return fmt.Errorf("account info name is different from specified name!\n")}

    if dbg {
        fmt.Printf("account info\n")
        fmt.Printf("Contacts:\n")
        for i :=0; i < len(acntInfo.Contacts); i++ {
            fmt.Printf("  %d: %s\n", i+1, acntInfo.Contacts[i])
        }
    }


    leAcnt := LEAcnt {
        AcntNam: LELibObj.Name,
        PrivKeyFilnam: privKeyFilnam,
        PubKeyFilnam: pubKeyFilnam,
        Contacts: acntInfo.Contacts,
    }

    if LELibObj.Prod {
		leAcnt.Prod = true
        leAcnt.LEUrl = LEProdUrl
    } else {
        leAcnt.Prod = false
        leAcnt.LEUrl = LETestUrl
    }

    akey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil { return fmt.Errorf("Generate Key: %v", err)}

    if dbg {fmt.Printf("dbg -- newClient: key generated!\n")}

    client := &acme.Client{
        Key: akey,
        DirectoryURL: leAcnt.LEUrl,
        }

    acntTpl:= acme.Account{
        Contact: leAcnt.Contacts,
    }

    acnt, err := client.Register(ctx, &acntTpl, acme.AcceptTOS)
    if err != nil { return fmt.Errorf("client.Register: %v", err)}

    if dbg {
        fmt.Printf("dbg -- Directory Url: %s\n", client.DirectoryURL)
        fmt.Printf("dbg -- success client and account created!\n")
//        PrintClient(client)
        PrintAcmeAccount(acnt)
    }

    privateKey := (client.Key).(*ecdsa.PrivateKey)

    publicKey := &ecdsa.PublicKey{}

    publicKey = &privateKey.PublicKey

    x509Encoded, err := x509.MarshalECPrivateKey(privateKey)
    if err != nil {return fmt.Errorf("x509.MarshalECPrivateKey: %v", err)}

    pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

    err = os.WriteFile(privKeyFilnam, pemEncoded, 0644)
    if err != nil {return fmt.Errorf("pem priv key write file: %v", err)}

    x509EncodedPub, err := x509.MarshalPKIXPublicKey(publicKey)
    if err != nil {return fmt.Errorf("x509.MarshalPKIXPublicKey: %v", err)}

    pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
    err = os.WriteFile(pubKeyFilnam, pemEncodedPub, 0644)
    if err != nil {return fmt.Errorf("pem pub key write file: %v", err)}

    leAcnt.Updated = time.Now()
    leAcnt.AcntId = string(client.KID)

    newAcntData, err := yaml.Marshal(&leAcnt)
    if err != nil {return fmt.Errorf("yaml Unmarshal account file: %v\n", err)}

    if err = os.WriteFile(acntFilnam, newAcntData, 0600); err != nil {
        return fmt.Errorf("Error writing key file %q: %v", acntFilnam, err)
    }

    return nil
}


// method that creates LE client object and verifies LE account
//func (certobj *CertObj) GetAcmeClientV2(ctx context.Context) (err error) {
func (LELibObj *LELib) GetLEAccount() (acnt *acme.Account, err error) {

//    var client acme.Client
    dbg := LELibObj.Dbg
	ctx := LELibObj.Ctx

	LEDir := os.Getenv("LEDir")
	if len(LEDir) == 0 {return nil, fmt.Errorf("cannot find LEDir!")}

// contruct LE account name
	acntFilnam :=""
	privKeyFilnam := ""
	pubKeyFilnam := ""

	if LELibObj.Prod {
    	acntFilnam =  LEDir + "/" + LELibObj.Name + "LEProd.yaml"
	    privKeyFilnam = LEDir + "/" + LELibObj.Name + "LEProdPriv.key"
    	pubKeyFilnam = LEDir + "/" + LELibObj.Name + "LEProdPub.key"
	} else {
    	acntFilnam =  LEDir + "/" + LELibObj.Name + "LETest.yaml"
	    privKeyFilnam = LEDir + "/" + LELibObj.Name + "LETestPriv.key"
    	pubKeyFilnam = LEDir + "/" + LELibObj.Name + "LETestPub.key"
	}

    if dbg {
		fmt.Println("*************** dbg info ***************")
		fmt.Printf("account file:    %s\n", acntFilnam)
//		fmt.Printf("account info:    %s\n", acntInfoFilnam)
        fmt.Printf("priv Key Filnam: %s\n", privKeyFilnam)
        fmt.Printf("pub Key Filnam:  %s\n", pubKeyFilnam)
		fmt.Println("*********** end dbg info ***************")
    }

    acntData, err := os.ReadFile(acntFilnam)
    if err != nil {return nil, fmt.Errorf("cannot read LE account file! %v", err)}

    leAcnt := LEAcnt{}

    err = yaml.Unmarshal(acntData, &leAcnt)
    if err != nil {return nil, fmt.Errorf("yaml Unmarshal account file: %v\n", err)}

    if dbg {PrintLEAcnt(&leAcnt)}

    pemEncoded, err := os.ReadFile(leAcnt.PrivKeyFilnam)
    if err != nil {return nil, fmt.Errorf("os.Read Priv Key: %v", err)}

    pemEncodedPub, err := os.ReadFile(leAcnt.PubKeyFilnam)
    if err != nil {return nil, fmt.Errorf("os.Read Pub Key: %v", err)}

    block, _ := pem.Decode([]byte(pemEncoded))
    x509Encoded := block.Bytes
    privateKey, err := x509.ParseECPrivateKey(x509Encoded)
    if err != nil {return nil, fmt.Errorf("x509.ParseECPivateKey: %v", err)}

    blockPub, _ := pem.Decode([]byte(pemEncodedPub))
    x509EncodedPub := blockPub.Bytes
    genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
    if err != nil {return nil, fmt.Errorf("x509.ParsePKIXKey: %v", err)}

    publicKey := genericPublicKey.(*ecdsa.PublicKey)
    privateKey.PublicKey = *publicKey

	client:= &acme.Client{}
    client.Key = privateKey
    client.DirectoryURL = leAcnt.LEUrl

    acnt, err = client.GetReg(ctx, "")
    if err != nil {return nil, fmt.Errorf("error -- LE GetReg account: %v\n", err)}

    if acnt.Status != "valid" {
        return acnt, fmt.Errorf("error -- LE acount is not valid. status: %s\n", acnt.Status)
    }

    return acnt, nil
}

func PrintLEAcnt(acnt *LEAcnt) {

    fmt.Printf("*************** LEAcnt *******************\n")
    fmt.Printf("Acnt Name:  %s\n", acnt.AcntNam)
    fmt.Printf("AcntId:     %s\n", acnt.AcntId)
    fmt.Printf("update:     %s\n", acnt.Updated.Format(time.RFC1123))
    fmt.Printf("LE Url:     %s\n", acnt.LEUrl)
    fmt.Printf("Prod:       %t\n", acnt.Prod)
    fmt.Printf("contacts:   %d\n", len(acnt.Contacts))
    for i:=0; i< len(acnt.Contacts); i++ {
        fmt.Printf("contact[%d]: %s\n", i+1, acnt.Contacts[i])
    }
    fmt.Printf("Public Key File:  %s\n", acnt.PubKeyFilnam)
    fmt.Printf("Private Key File: %s\n", acnt.PrivKeyFilnam)
    fmt.Printf("*************** End LEAcnt ****************\n")
}

func (leObj *LELib) PrintLELibObj() {

    fmt.Printf("*************** LELibObj *******************\n")
	fmt.Printf("Name:  %s\n", leObj.Name)
	fmt.Printf("LEDir: %s\n", leObj.LEDir)
	fmt.Printf("Prod:  %t\n", leObj.Prod)
	fmt.Printf("Dbg:   %t\n", leObj.Dbg)
    fmt.Printf("************* End LELibObj *****************\n")
}

func PrintAcmeAccount(acnt *acme.Account) {

    fmt.Println("***************** Acme Account ******************")
    fmt.Printf("URI:    %s\n", acnt.URI)
    fmt.Printf("Status: %s\n", acnt.Status)
    fmt.Printf("Contacts [%d]:\n", len((*acnt).Contact))
    for i:=0; i< len((*acnt).Contact); i++ {
        fmt.Printf("Contact[%d]: %s\n", i, (*acnt).Contact[i])
    }
    fmt.Printf("OrdersURL:   %s\n", acnt.OrdersURL)
    fmt.Println (" *** non RFC 8588 terms:  ***")
    fmt.Printf("  AgreedTerms: %s\n", acnt.AgreedTerms)
    fmt.Printf("  Authz: %s\n", acnt.Authz)
    fmt.Println("***************** End Acme Account ******************")
}
