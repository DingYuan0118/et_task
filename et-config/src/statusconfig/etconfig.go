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
	StatusSuccess             = 0    // |0|成功|
	StatusServerError         = 1000 // |1000|服务器错误|
	StatusInvalidParams       = 1001 // |1001|非法参数|
	StatusNotFound            = 1002 // |1002|Not found|
	StatusLoginFailed         = 2001 // |2001|登录失败|
	StatusTokenInvalid        = 2002 // |2002|Token 失效，重新登录|
	StatusNoToken             = 2003 // |2003|请求头中无token，重新登录|
	StatusQueryFaild          = 3001 // |3001|查询失败|
	StatusUpdateNicknameFaild = 3002 // |3002|更新昵称失败|
	StatusUploadPicFailed     = 3003 // |3003|上传头像失败|
	StatusThirdPackageErr     = 4000 // |4000|第三方库返回错误|
)

var ErrMsg = map[int]string{
	StatusSuccess:             "Success",                                                    // |0|成功|
	StatusServerError:         "Server Error",                                               // |1000|服务器错误|
	StatusInvalidParams:       "Invalid Params",                                             // |1001|非法参数|
	StatusNotFound:            "Not Found",                                                  // |1002|Not found|
	StatusLoginFailed:         "Login Failed",                                               // |2001|登录失败|
	StatusTokenInvalid:        "Token Invalid",                                              // |2002|Token 失效，重新登录|
	StatusNoToken:             "Token no exist",                                             // |2003|请求头中无token，重新登录|
	StatusQueryFaild:          "Query Faild",                                                // |3001|查询失败|
	StatusUpdateNicknameFaild: "Update Nickname Faild: Nickname too long or user not exits", // |3002|更新昵称失败|
	StatusUploadPicFailed:     "Upload Picture Failed",                                      // |3003|上传头像失败|
}

// 图像存储路径
const ImageFolder string = "/Users/yuan.ding/Desktop/code/entry_task/userimages/"