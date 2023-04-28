package types

const (
	OneAnswerQuestionCode  int8 = 1
	MultAnswerQuestionCode int8 = 2
	TrueFalseQuestionCode  int8 = 3
	MatchQuestionCode      int8 = 4
	OrderQuestionCode      int8 = 5
	AdminRoleCode          int8 = 0
	TeacherRoleCode        int8 = 1
	StudentRoleCode        int8 = 2
)

type Question struct {
	Text             string
	AnswerOptions    string
	CorrectAnswer    string
	QuestionTypeCode int8
	Score            int64
	Terms            []string
}

type Test struct {
	Title      string
	Subject    string
	BeginTime  string
	EndTime    string
	IsAdaptive bool
	Groups     []string
	Terms      []string
	Questions  []Question
}

type User struct {
	Name       string
	Surname    string
	LastName   string
	Login      string
	Password   string
	RoleCode   int8
	TelegramID string
}
