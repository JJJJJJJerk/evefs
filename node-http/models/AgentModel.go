package models

type Agent struct {
	Id        uint `gorm:"primary_key"`
	AgentName string
	Phone     string
	ParentId  int
}

func (Agent) TableName() string {
	return "a_agent"
}

func GetAllAgents() (agents []Agent) {
	DB.Find(&agents)
	return
}
