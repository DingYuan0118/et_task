package statusconfig

// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.
var HmacSampleSecret = []byte("mHpdHzQtEWQw7ntdpoNe")

// gRPC 设置
const (
	DefaultName = "yuan"
)

// var (
// 	addr = flag.String("addr", "localhost:50051", "the address to connect to")
// 	name = flag.String("name", defaultName, "Name to greet")
// )

// 状态码设置
const (
	StatusSuccess                   = 0    // |0|成功|
	StatusServerError               = 1000 // |1000|服务器错误|
	StatusInvalidParams             = 1001 // |1001|非法参数|
	StatusNotFound                  = 1002 // |1002|Not found|
	StatusLoginFailedPasswordWrong  = 2001 // |2001|登录失败，密码错误|
	StatusLoginFailedNoUser         = 2002 // |2002|登录失败，用户不存在|
	StatusTokenInvalid              = 2003 // |2003|Token 失效，重新登录|
	StatusNoToken                   = 2004 // |2004|请求头中无token，重新登录|
	StatusQueryFaild                = 3001 // |3001|查询失败|
	StatusNicknameTooLong           = 4001 // |4001|更新昵称失败,昵称过长|
	StatusUpdateNicknameFaildNoUser = 4002 // |4002|更新昵称失败,用户不存在|
	StatusUploadPicFormatWrong      = 5001 // |5001|上传头像失败,格式错误|
	StatusUploadPicTooLarge         = 5002 // |5002|上传头像失败,文件过大|
	StatusUploadPicFailedNouser     = 5003 // |5003|上传头像失败,用户不存在|
)

var ErrMsg = map[int]string{
	StatusSuccess:                   "Success. ",
	StatusServerError:               "Server Error. ",
	StatusInvalidParams:             "Invalid Params. ",
	StatusNotFound:                  "Not Found. ",
	StatusLoginFailedPasswordWrong:  "Login Failed. Wrong Password",
	StatusLoginFailedNoUser:         "Login Failed. User not exist",
	StatusTokenInvalid:              "Token Invalid. Please Login again",
	StatusNoToken:                   "Token not exist. ",
	StatusQueryFaild:                "Query Faild. User not exist",
	StatusNicknameTooLong:           "Update Nickname Faild. Nickname Too Long, should less than 64",
	StatusUpdateNicknameFaildNoUser: "Update Nickname Faild. User not exist",
	StatusUploadPicFormatWrong:      "Upload Picture Failed. file format error, support [png, jpeg, bmp]",
	StatusUploadPicTooLarge:         "Upload Picture Failed. file too large, should less than 3MB",
	StatusUploadPicFailedNouser:     "Upload Picture Failed. User not exist",
}

// 图像存储路径
const ImageFolder string = "/Users/yuan.ding/Desktop/code/entry_task/userimages/"
