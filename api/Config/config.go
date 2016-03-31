/* CONFIG.GO SHOULD BE GITIGNORED -- DB & PORT INFO */
package config

var Port string = ":8000"
var DB string = "mongodb://localhost:27017/dartboard"
var Secret string = ""
var JwtSecret string = ""
var Email = &mail{"emailname", "password", "smtp.gmail.com", 587}
type mail struct {
    Username    string
    Password    string
    EmailServer string
    Port        int
}
