package handlers

import (
	"log"
	"net/http"
	"resume-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OnlineResumeHandler 在线简历处理器
type OnlineResumeHandler struct {
	DB *gorm.DB
}

// NewOnlineResumeHandler 创建在线简历处理器
func NewOnlineResumeHandler(db *gorm.DB) *OnlineResumeHandler {
	// 自动迁移在线简历表
	if err := db.AutoMigrate(&models.OnlineResume{}); err != nil {
		log.Printf("Warning: Failed to migrate OnlineResume table: %v", err)
	}
	return &OnlineResumeHandler{DB: db}
}

// GetOnlineResume 获取当前用户的在线简历
// GET /resumes/online
// Requirements: 4.1 - 显示当前简历信息
// 如果用户没有在线简历，则创建一个空的简历记录
func (h *OnlineResumeHandler) GetOnlineResume(c *gin.Context) {
	// 从 JWT 中获取用户 ID
	jwtUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权，请先登录",
		})
		return
	}

	userID, ok := jwtUserID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的用户ID",
		})
		return
	}

	var resume models.OnlineResume

	// 查找用户的在线简历
	result := h.DB.Where("user_id = ?", userID).First(&resume)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 如果不存在，创建一个空的在线简历
			resume = models.OnlineResume{
				UserID:         userID,
				WorkExperience: models.WorkExperienceList{},
				Education:      models.EducationList{},
				Skills:         models.SkillList{},
			}

			if err := h.DB.Create(&resume).Error; err != nil {
				log.Printf("[在线简历] 创建空简历失败: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "创建简历失败",
				})
				return
			}

			log.Printf("[在线简历] ✓ 为用户 %d 创建了空的在线简历", userID)
		} else {
			log.Printf("[在线简历] 查询失败: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取简历失败",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    resume.ToResponse(),
	})
}

// SaveOnlineResume 保存在线简历
// PUT /resumes/online
// Requirements: 4.3 - 持久化更新的信息到数据库
// Requirements: 4.6 - 验证必填字段（name, phone, email）
func (h *OnlineResumeHandler) SaveOnlineResume(c *gin.Context) {
	// 从 JWT 中获取用户 ID
	jwtUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权，请先登录",
		})
		return
	}

	userID, ok := jwtUserID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的用户ID",
		})
		return
	}

	// 解析请求体
	var req models.OnlineResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 查找或创建用户的在线简历
	var resume models.OnlineResume
	result := h.DB.Where("user_id = ?", userID).First(&resume)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 创建新的在线简历
			resume = models.OnlineResume{
				UserID: userID,
			}
		} else {
			log.Printf("[在线简历] 查询失败: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "保存简历失败",
			})
			return
		}
	}

	// 从请求更新简历数据
	resume.FromRequest(req)

	// Requirements 4.6: 验证必填字段（name, phone, email）
	if err := resume.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 保存到数据库
	var saveErr error
	if resume.ID == 0 {
		// 新建
		saveErr = h.DB.Create(&resume).Error
	} else {
		// 更新
		saveErr = h.DB.Save(&resume).Error
	}

	if saveErr != nil {
		log.Printf("[在线简历] 保存失败: %v", saveErr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存简历失败",
		})
		return
	}

	log.Printf("[在线简历] ✓ 用户 %d 的在线简历已保存", userID)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "简历保存成功",
		"data":    resume.ToResponse(),
	})
}
