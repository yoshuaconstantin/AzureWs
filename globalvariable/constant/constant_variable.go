package constant

const BcryptCost = 14

// Secret key should be hidden!!
var Salt = []byte("AzureKey")

const CreateAccountPost = "/api/create-account"
const UserDeleteAccount = "/api/user/delete-account"
const UserUpdatePassword = "/api/user/update-password"
const UserLogin = "/api/login"
const UserLogout = "/api/logout"
const HomeDashboards = "/api/home/dashboards"
const HomeUpdateDashboards = "/api/home/update/dashboards/data"
const HomeUserProfileImage = "/api/home/user/profile/image"
const HomeUserProfile = "/api/home/user/profile"
const HomeUserFeedback = "/api/home/user/feedback"
const CommunityPost = "/api/community/post"
const CommunityPostLike = "/api/community/post/like"
const CommunityPostComment = "/api/community/post/comment"
const CommunityPostCommentGet = "/api/community/post/comment/get-comment"
const RefreshToken = "/api/token-refresh"
