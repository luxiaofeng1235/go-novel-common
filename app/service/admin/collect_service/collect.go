package collect_service

import (
	"encoding/json"
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"strings"
)

func GetCollectById(id int64) (collect *models.McCollect, err error) {
	err = global.DB.Model(models.McCollect{}).Where("id", id).First(&collect).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCollectResById(id int64) (collectRes *models.GetCollectRes, err error) {
	var collect *models.McCollect
	collect, err = GetCollectById(id)
	if err != nil {
		return
	}
	categoryArr := []*models.CategoryReg{}
	if collect.Categorys != "" {
		err = json.Unmarshal([]byte(collect.Categorys), &categoryArr)
		if err != nil {
			return
		}
	}
	listPageArr := []*models.ListPageReg{}
	if collect.ListPageReg != "" {
		err = json.Unmarshal([]byte(collect.ListPageReg), &listPageArr)
		if err != nil {
			return
		}
	}
	descReplaceArr := []*models.TextReplaceReg{}
	if collect.DescReplaceReg != "" {
		log.Println(collect.DescReplaceReg)
		err = json.Unmarshal([]byte(collect.DescReplaceReg), &descReplaceArr)
		if err != nil {
			log.Println("err11", err.Error())
			return
		}
	}

	textReplaceArr := []*models.TextReplaceReg{}
	if collect.TextReplaceReg != "" {
		err = json.Unmarshal([]byte(collect.TextReplaceReg), &textReplaceArr)
		if err != nil {
			return
		}
	}

	collectRes = &models.GetCollectRes{
		Id:                collect.Id,
		Title:             collect.Title,
		Link:              collect.Link,
		Charset:           collect.Charset,
		UrlComplete:       collect.UrlComplete,
		UrlReverse:        collect.UrlReverse,
		PicLocal:          collect.PicLocal,
		CategoryArr:       categoryArr,
		CategoryWay:       collect.CategoryWay,
		CategoryFixed:     collect.CategoryFixed,
		ListPageArr:       listPageArr,
		ListSectionReg:    collect.ListSectionReg,
		ListUrlReg:        collect.ListUrlReg,
		ChapterSectionReg: collect.ChapterSectionReg,
		ChapterUrlReg:     collect.ChapterUrlReg,
		ChapterTextReg:    collect.ChapterTextReg,
		CategoryNameReg:   collect.CategoryNameReg,
		BookNameReg:       collect.BookNameReg,
		DescReg:           collect.DescReg,
		PicReg:            collect.PicReg,
		AuthorReg:         collect.AuthorReg,
		Status:            collect.Status,
		SerializeReg:      collect.SerializeReg,
		TagNameReg:        collect.TagNameReg,
		UpdateReg:         collect.UpdateReg,
		DescReplaceArr:    descReplaceArr,
		TextReplaceArr:    textReplaceArr,
	}
	return
}

func CollectListSearch(req *models.CollectListReq) (list []*models.McCollect, total int64, err error) {
	db := global.DB.Model(&models.McCollect{}).Order("id desc")

	title := strings.TrimSpace(req.Title)
	if title != "" {
		db = db.Where("title = ?", title)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

func CreateCollect(req *models.CreateCollectReq) (InsertId int64, err error) {
	title := strings.TrimSpace(req.Title)
	link := strings.TrimSpace(req.Link)
	charset := strings.TrimSpace(req.Charset)
	listSectionReg := strings.TrimSpace(req.ListSectionReg)
	listUrlReg := strings.TrimSpace(req.ListUrlReg)
	chapterSectionReg := strings.TrimSpace(req.ChapterSectionReg)
	chapterUrlReg := strings.TrimSpace(req.ChapterUrlReg)
	chapterTextReg := strings.TrimSpace(req.ChapterTextReg)
	categoryNameReg := strings.TrimSpace(req.CategoryNameReg)
	bookNameReg := strings.TrimSpace(req.BookNameReg)
	picReg := strings.TrimSpace(req.PicReg)
	descReg := strings.TrimSpace(req.DescReg)
	authorReg := strings.TrimSpace(req.AuthorReg)
	serializeReg := strings.TrimSpace(req.SerializeReg)
	updateReg := strings.TrimSpace(req.UpdateReg)
	tagNameReg := strings.TrimSpace(req.TagNameReg)

	listPageArr := req.ListPageArr
	categoryArr := req.CategoryArr
	descReplaceArr := req.DescReplaceArr
	textReplaceArr := req.TextReplaceArr

	urlComplete := req.UrlComplete
	urlReverse := req.UrlReverse
	picLocal := req.PicLocal
	categoryWay := req.CategoryWay
	categoryFixed := req.CategoryFixed
	status := req.Status

	if title == "" {
		err = fmt.Errorf("%v", "规则名称不能为空")
		return
	}

	if link == "" {
		err = fmt.Errorf("%v", "采集链接不能为空")
		return
	}

	if charset == "" {
		err = fmt.Errorf("%v", "网站编码不能为空")
		return
	}

	if len(listPageArr) <= 0 {
		err = fmt.Errorf("%v", "分页列表地址不能为空")
		return
	}
	if listSectionReg == "" {
		err = fmt.Errorf("%v", "列表区间正则不能为空")
		return
	}
	if listUrlReg == "" {
		err = fmt.Errorf("%v", "小说详情链接正则不能为空")
		return
	}
	if categoryNameReg == "" {
		err = fmt.Errorf("%v", "栏目正则不能为空")
		return
	}
	if categoryWay <= 0 && len(categoryArr) <= 0 {
		err = fmt.Errorf("%v", "栏目转换不能为空")
		return
	}

	if categoryWay > 0 && categoryFixed <= 0 {
		err = fmt.Errorf("%v", "小说分类不能为空")
		return
	}

	if bookNameReg == "" {
		err = fmt.Errorf("%v", "小说名称正则正则不能为空")
		return
	}

	if descReg == "" {
		err = fmt.Errorf("%v", "小说简介正则不能为空")
		return
	}

	if picReg == "" {
		err = fmt.Errorf("%v", "小说图片正则不能为空")
		return
	}

	if authorReg == "" {
		err = fmt.Errorf("%v", "小说作者正则不能为空")
		return
	}

	if serializeReg == "" {
		err = fmt.Errorf("%v", "小说连载状态正则不能为空")
		return
	}

	if updateReg == "" {
		err = fmt.Errorf("%v", "小说更新时间正则不能为空")
		return
	}

	if tagNameReg == "" {
		err = fmt.Errorf("%v", "小说标签正则不能为空")
		return
	}

	if chapterSectionReg == "" {
		err = fmt.Errorf("%v", "章节区间正则不能为空")
		return
	}

	if chapterUrlReg == "" {
		err = fmt.Errorf("%v", "章节链接正则不能为空")
		return
	}
	if chapterTextReg == "" {
		err = fmt.Errorf("%v", "小说内容正则不能为空")
		return
	}
	collect := models.McCollect{
		Title:             title,
		Link:              link,
		Charset:           charset,
		Status:            status,
		UrlComplete:       urlComplete,
		UrlReverse:        urlReverse,
		PicLocal:          picLocal,
		ListPageReg:       utils.JSONString(listPageArr),
		ListSectionReg:    listSectionReg,
		ListUrlReg:        listUrlReg,
		CategoryNameReg:   categoryNameReg,
		CategoryWay:       categoryWay,
		Categorys:         utils.JSONString(categoryArr),
		CategoryFixed:     categoryFixed,
		BookNameReg:       bookNameReg,
		AuthorReg:         authorReg,
		SerializeReg:      serializeReg,
		PicReg:            picReg,
		UpdateReg:         updateReg,
		TagNameReg:        tagNameReg,
		DescReg:           descReg,
		DescReplaceReg:    utils.JSONString(descReplaceArr),
		ChapterSectionReg: chapterSectionReg,
		ChapterUrlReg:     chapterUrlReg,
		ChapterTextReg:    chapterTextReg,
		TextReplaceReg:    utils.JSONString(textReplaceArr),
		Addtime:           utils.GetUnix(),
	}

	if err = global.DB.Create(&collect).Error; err != nil {
		return 0, err
	}
	return collect.Id, nil
}

func UpdateCollect(req *models.UpdateCollectReq) (res bool, err error) {
	id := req.CollectId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不正确")
		return
	}
	title := strings.TrimSpace(req.Title)
	link := strings.TrimSpace(req.Link)

	listSectionReg := strings.TrimSpace(req.ListSectionReg)
	listUrlReg := strings.TrimSpace(req.ListUrlReg)
	chapterSectionReg := strings.TrimSpace(req.ChapterSectionReg)
	chapterUrlReg := strings.TrimSpace(req.ChapterUrlReg)
	chapterTextReg := strings.TrimSpace(req.ChapterTextReg)
	categoryNameReg := strings.TrimSpace(req.CategoryNameReg)
	bookNameReg := strings.TrimSpace(req.BookNameReg)
	descReg := strings.TrimSpace(req.DescReg)
	picReg := strings.TrimSpace(req.PicReg)
	authorReg := strings.TrimSpace(req.AuthorReg)
	serializeReg := strings.TrimSpace(req.SerializeReg)
	updateReg := strings.TrimSpace(req.UpdateReg)
	tagNameReg := strings.TrimSpace(req.TagNameReg)
	listPageArr := req.ListPageArr
	categoryArr := req.CategoryArr
	descReplaceArr := req.DescReplaceArr
	textReplaceArr := req.TextReplaceArr

	charset := strings.TrimSpace(req.Charset)
	urlComplete := req.UrlComplete
	urlReverse := req.UrlReverse
	picLocal := req.PicLocal
	categoryFixed := req.CategoryFixed
	status := req.Status
	categoryWay := req.CategoryWay

	var mapData = make(map[string]interface{})
	if title != "" {
		mapData["title"] = title
	}
	if link != "" {
		mapData["link"] = link
	}
	if charset != "" {
		mapData["charset"] = charset
	}

	if listSectionReg != "" {
		mapData["list_section_reg"] = listSectionReg
	}
	if listUrlReg != "" {
		mapData["list_url_reg"] = listUrlReg
	}

	if chapterSectionReg != "" {
		mapData["chapter_section_reg"] = chapterSectionReg
	}
	if chapterUrlReg != "" {
		mapData["chapter_url_reg"] = chapterUrlReg
	}
	if chapterTextReg != "" {
		mapData["chapter_text_reg"] = chapterTextReg
	}
	if categoryNameReg != "" {
		mapData["category_name_reg"] = categoryNameReg
	}
	if bookNameReg != "" {
		mapData["book_name_reg"] = bookNameReg
	}
	if descReg != "" {
		mapData["desc_reg"] = descReg
	}
	if picReg != "" {
		mapData["pic_reg"] = picReg
	}
	if authorReg != "" {
		mapData["author_reg"] = authorReg
	}
	if serializeReg != "" {
		mapData["serialize_reg"] = serializeReg
	}
	if updateReg != "" {
		mapData["update_reg"] = updateReg
	}
	if tagNameReg != "" {
		mapData["tag_name_reg"] = tagNameReg
	}
	if len(listPageArr) > 0 {
		mapData["list_page_reg"] = utils.JSONString(listPageArr)
	}
	if len(categoryArr) > 0 {
		mapData["categorys"] = utils.JSONString(categoryArr)
	}
	if len(descReplaceArr) > 0 {
		mapData["desc_replace_reg"] = utils.JSONString(descReplaceArr)
	}
	if len(textReplaceArr) > 0 {
		mapData["text_replace_reg"] = utils.JSONString(textReplaceArr)
	}
	mapData["url_complete"] = urlComplete
	mapData["url_reverse"] = urlReverse
	mapData["pic_local"] = picLocal
	mapData["category_fixed"] = categoryFixed
	mapData["status"] = status
	mapData["category_way"] = categoryWay
	mapData["uptime"] = utils.GetUnix()
	if err = global.DB.Model(models.McCollect{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DeleteCollect(req *models.DeleteCollectReq) (res bool, err error) {
	ids := req.CollectIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.McCollect{}).Error
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return
	}
	return true, nil
}
