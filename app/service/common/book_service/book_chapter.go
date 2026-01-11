package book_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
)

func GetChapterTable(bookId int64) (string, error) {
	var tabPrefix = "mc_"
	table := fmt.Sprintf("%vbook_chapter_%d", tabPrefix, bookId)

	// 数据分表是否存在
	if !global.DB.Migrator().HasTable(table) {
		// 表不存在，创建表
		err := createChapterTable(table)
		if err != nil {
			global.Sqllog.Errorf("table=%v 生成章节表失败", table)
			return "", err
		}
	}

	return table, nil
}

func createChapterTable(table string) (err error) {
	// 创建表
	err = global.DB.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INT(11) NOT NULL AUTO_INCREMENT,
			sort INT(11) DEFAULT '0' COMMENT '排序ID',
		    chapter_link VARCHAR(200) DEFAULT '' COMMENT '章节链接',
			chapter_name VARCHAR(200) DEFAULT '' COMMENT '章节标题',
			vip TINYINT(1) DEFAULT '0' COMMENT 'VIP阅读，0否1是',
			cion INT(11) DEFAULT '0' COMMENT '章节需要金币',
			text_num INT(11) DEFAULT '0' COMMENT '章节字数',
			addtime INT(11) DEFAULT '0' COMMENT '入库时间',
			PRIMARY KEY (id),
			KEY sort (sort) USING BTREE
		) ENGINE=MyISAM AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='小说章节';
	`, table)).Error
	return err
}

func GetChapterNames(tableName string) (chapterNames []string) {
	var err error
	err = global.DB.Table(tableName).Select("chapter_name").Order(fmt.Sprintf("sort asc")).Find(&chapterNames).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetSortFirst(tableName string) (sort int) {
	var err error
	err = global.DB.Table(tableName).Order("sort asc").Select("sort").Limit(1).Scan(&sort).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetChapterIdLast(tableName string) (chapter *models.McBookChapter) {
	var err error
	err = global.DB.Table(tableName).Order("sort DESC").Limit(1).Last(&chapter).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetBookNewChapterName(bookId int64) (newChapterName string) {
	//章节表
	chapterTable, err := GetChapterTable(bookId)
	if err != nil {
		err = fmt.Errorf("%v", "获取章节失败")
		return
	}
	err = global.DB.Table(chapterTable).Order("sort desc,id desc").Select("chapter_name").Limit(1).Scan(&newChapterName).Error
	if err != nil {
		return
	}
	return
}

func GetChapterList(tableName string, sort string) (chapters []*models.McBookChapter, err error) {
	db := global.DB.Table(tableName).Order(fmt.Sprintf("sort %v", sort))
	err = db.Find(&chapters).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
