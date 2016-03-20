/* CONFIG.GO SHOULD BE GITIGNORED -- DB & PORT INFO */
package config

var Port string = ":8000"
var DB string = ""
var Secret string = ""
var JwtSecret string = ""
var Email = &emailUser{"emailname", "password", "smtp.gmail.com", 587}
type emailUser struct {
    Username    string
    Password    string
    EmailServer string
    Port        int
}
