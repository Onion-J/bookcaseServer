package api

import (
	"BookcaseServer/common"
	"BookcaseServer/dto"
	"BookcaseServer/model"
	"BookcaseServer/response"
	"github.com/gin-gonic/gin"
	"sort"
)

type StudentData struct {
	StudentId   string          `json:"studentId"`
	AccountList []model.Student `json:"accountList"`
}

// GetEnrollmentYear 获取入学年份
func GetEnrollmentYear(c *gin.Context) {
	var enrollmentYear []string

	db := common.GetDB()

	// 查询
	db.Raw("SELECT SUBSTRING(student_id, 2, 2)  FROM students").Scan(&enrollmentYear)

	// 先排序，排序后两两比较,相同则去除
	sort.Strings(enrollmentYear)
	for i := len(enrollmentYear) - 1; i > 0; i-- {
		if enrollmentYear[i] == enrollmentYear[i-1] {
			enrollmentYear = append(enrollmentYear[:i], enrollmentYear[i+1:]...)
		}
	}

	// 返回结果
	response.Success(c, response.Selected, gin.H{"enrollmentYearList": enrollmentYear})
}

// CreateStudentAccount 创建学生账户
func CreateStudentAccount(c *gin.Context) {
	var studentData StudentData

	// 参数绑定
	if err := c.ShouldBind(&studentData); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if len(studentData.AccountList) == 0 {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 判断学生账户是否存在 存在则跳过 不存在则创建
	for _, student := range studentData.AccountList {
		if row := db.Where("student_id = ?", student.StudentId).First(&model.Student{}).RowsAffected; row == 0 {
			if err := db.Create(&student).Error; err != nil {
				response.Failed(c, response.NotCreated)
				return
			}
		}
	}

	// 返回结果
	response.Success(c, response.Created, nil)
}

// SelectStudentAccount 查询学生账户
func SelectStudentAccount(c *gin.Context) {
	var studentData StudentData

	// 参数绑定
	if err := c.ShouldBind(&studentData); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if studentData.StudentId == "" {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 查询
	var students []model.Student
	if len(studentData.StudentId) == 2 {
		// 学院
		db.Raw("SELECT * FROM students WHERE SUBSTRING(student_id, 4, 2) = ?", studentData.StudentId).Scan(&students)
	} else if len(studentData.StudentId) == 3 {
		// 年份
		db.Raw("SELECT * FROM students WHERE SUBSTRING(student_id, 1, 3) = ?", studentData.StudentId).Scan(&students)
	} else if len(studentData.StudentId) == 4 {
		// 专业
		db.Raw("SELECT * FROM students WHERE SUBSTRING(student_id, 4, 4) = ?", studentData.StudentId).Scan(&students)
	} else if len(studentData.StudentId) == 5 {
		// 年份 + 学院
		db.Raw("SELECT * FROM students WHERE SUBSTRING(student_id, 1, 5) = ?", studentData.StudentId).Scan(&students)
	} else if len(studentData.StudentId) == 7 {
		// 年份 + 专业
		db.Raw("SELECT * FROM students WHERE SUBSTRING(student_id, 1, 7) = ?", studentData.StudentId).Scan(&students)
	} else if len(studentData.StudentId) == 12 {
		db.Where("student_id = ?", studentData.StudentId).Find(&students)
	} else {
		response.Failed(c, response.DataError)
		return
	}

	// 类型转换
	var studentsDto = make([]dto.StudentDto, len(students))
	for i, v := range students {
		studentsDto[i] = dto.ToStudentDto(v)
	}

	// 返回结果
	response.Success(c, response.Selected, gin.H{"students": studentsDto})

}

// DeleteStudentAccount 删除学生账户
func DeleteStudentAccount(c *gin.Context) {
	var studentData StudentData

	// 参数绑定
	if err := c.ShouldBind(&studentData); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if len(studentData.AccountList) == 0 {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 判断学生账户是否存在 不存在则跳过 存在则删除
	for _, student := range studentData.AccountList {
		if row := db.Where("student_id = ?", student.StudentId).First(&model.Student{}).RowsAffected; row != 0 {
			if err := db.Delete(&student).Error; err != nil {
				response.Failed(c, response.NotDeleted)
				return
			}
		}
	}

	// 返回结果
	response.Success(c, response.Deleted, nil)
}
