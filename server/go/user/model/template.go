package model

import (
	templatei "github.com/hopeio/zeta/utils/definition/template"
)

const (
	ActionActiveContent  = `ActionActiveContent`
	ActionActiveTemplate = `{{define "ActionActiveContent"}}<p><b>亲爱的{{.UserName}}:</b></p>
			<p>我们收到您在 {{.SiteName}} 的注册信息, 请点击下面的链接, 或粘贴到浏览器地址栏来激活帐号.</p>
			<a href="{{.ActionURL}}">{{.ActionURL}}</a>
			<p>如果您没有在 {{.SiteName}} 填写过注册信息, 说明有人滥用了您的邮箱, 请删除此邮件, 我们对给您造成的打扰感到抱歉.</p>
			<p>{{.SiteName}} 谨上.</p>{{end}}`
	ActionRestPasswordContent  = `ActionRestPasswordContent`
	ActionRestPasswordTemplate = `{{define "ActionRestPasswordContent"}}<p><b>亲爱的{{.UserName}}:</b></p>
			<p>你的密码重设要求已经得到验证。请点击以下链接, 或粘贴到浏览器地址栏来设置新的密码: </p>
			<a href="{{.ActionURL}}">{{.ActionURL}}</a>
			<p>感谢你对 {{.SiteName}} 的支持，希望你在 {{.SiteName}} 的体验有益且愉快。</p>
			<p>(这是一封自动产生的email，请勿回复。)</p>{{end}}`
)

func init() {
	templatei.Parse(ActionActiveTemplate)
	templatei.Parse(ActionRestPasswordTemplate)
}
