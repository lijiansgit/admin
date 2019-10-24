package ldap

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/lijiansgit/admin/config"
	log "github.com/lijiansgit/go/libs/log4go"
	"gopkg.in/ldap.v3"
)

var (
	LDAP *LDAPService
)

type LDAPService struct {
	Conn *ldap.Conn
}

func NewLDAPService() (ldapService *LDAPService, err error) {
	ldapService = new(LDAPService)
	// todo tcp keep alive
	// conn, err := ldap.Dial("tcp", config.Conf.LDAP.Addr)
	// if err != nil {
	// 	return ldapService, err
	// }

	// conn.SetTimeout(10 * time.Second)
	// // Reconnect with TLS
	// err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	// if err != nil {
	// 	return ldapService, err
	// }

	// err = conn.Bind(config.Conf.LDAP.Username, config.Conf.LDAP.Password)
	// if err != nil {
	// 	return ldapService, err
	// }

	// log.Info("Connect ldap success: %s", config.Conf.LDAP.Username)

	// ldapService.Conn = conn
	return ldapService, nil
}

func (l *LDAPService) CreateConn() (err error) {
	l.Conn, err = ldap.Dial("tcp", config.Conf.LDAP.Addr)
	if err != nil {
		return err
	}

	l.Conn.SetTimeout(10 * time.Second)
	// Reconnect with TLS
	err = l.Conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return err
	}

	err = l.Conn.Bind(config.Conf.LDAP.Username, config.Conf.LDAP.Password)
	if err != nil {
		return err
	}

	log.Debug("Connect ldap success: %s", config.Conf.LDAP.Username)
	return nil
}

func (l *LDAPService) Login(username, password string) (bool, error) {
	return true, nil // test skip ldap
	if err := l.CreateConn(); err != nil {
		return false, err
	}

	defer l.Conn.Close()

	searchRequest := ldap.NewSearchRequest(
		config.Conf.LDAP.RootDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(sAMAccountName=%s)", username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Conn.Search(searchRequest)
	if err != nil {
		return false, err
	}

	if len(sr.Entries) != 1 {
		return false, fmt.Errorf(
			"User: %s does not exist or too many entries returned", username)
	}

	userdn := sr.Entries[0].DN
	// Bind as the user to verify their password
	err = l.Conn.Bind(userdn, password)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (l *LDAPService) GetUsers() (users map[string]string, err error) {
	if err := l.CreateConn(); err != nil {
		return users, err
	}

	defer l.Conn.Close()

	searchRequest := ldap.NewSearchRequest(
		config.Conf.LDAP.RootDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=organizationalPerson))",
		[]string{"sAMAccountName", "mail", "dn"},
		// []string{"dn"},
		nil,
	)

	sr, err := l.Conn.Search(searchRequest)
	if err != nil {
		return users, err
	}

	users = make(map[string]string)
	for _, entry := range sr.Entries {
		// fmt.Println(entry.DN) CN=周青松,OU=tech,DC=tjj,DC=work
		name := entry.GetAttributeValue("sAMAccountName")
		email := entry.GetAttributeValue("mail")
		users[name] = email
	}

	return users, nil
}

func Init() (err error) {
	LDAP, err = NewLDAPService()
	if err != nil {
		return err
	}

	return nil
}
