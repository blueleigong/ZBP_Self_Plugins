package twodwife

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "math/rand"
    "time"
    "strconv"
    "hash/fnv"
    
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)



// 定义属性
var lpType = []string{"幼女", "萝莉", "万年萝莉", "合法萝莉", "萝莉老太婆", "乙女", "御姐", "非法御姐", "可萝可御", "软妹", "熟女", "人妻", "幼妻", "BBA", "伪娘", "秀吉", "伪伪娘", "假小子"}
var lpColorEyes = []string{"点缀着星空的夜之色", "棕红色", "湖蓝色", "蓝色", "蓝白色", "金色", "丹霞橙", "纯黑色", "棕红色", "淡蓝色", "稻草色", " 棕色", "灰色", "蓝黑色", "罗兰紫", "夜光黑", "白色", "淡金色", "青山黛", "豆芽绿色", "橘黄色", "深红色", "湖蓝色", "西瓜红", "星河银", "紫黄渐变", "橙红色", "浅紫色", "深绿色", "宝石蓝", "蜜糖色", "蓝绿色", "黑白闪烁", "橙色", "碎金色", "鹅黄色", "浅橙色", "枫叶色", "翡冷翠", "紫色", "蓝红渐变", "可怜的原谅色", "海藻绿", "深蓝色", "墨绿色", "灰色", "草绿色", "花色", "灰绿色", "奶白色", "银色", "黑灰色", "黄绿色", "粉白色", "嫩绿色", "奶金色", "浅绿色", "绛紫色", "浅紫色", "豆绿色", "黄琥珀色", "红色", "点缀着星空的夜之色", "棕红色", "湖蓝色", "蓝色", "蓝白色", "金色", "丹霞橙", "纯黑色", "棕红色", "淡蓝色", "稻草色", " 棕色", "灰色", "蓝黑色", "罗兰紫", "夜光黑", "白色", "淡金色", "青山黛", "豆芽绿色", "橘黄色", "深红色", "湖蓝色", "西瓜红", "星河银", "紫黄渐变", "橙红色", "浅紫色", "深绿色", "宝石蓝", "蜜糖色", "蓝绿色", "黑白闪烁", "橙色", "碎金色", "鹅黄色", "浅橙色", "枫叶色", "翡冷翠", "紫色", "蓝红渐变", "可怜的原谅色", "海藻绿", "深蓝色", "墨绿色", "灰色", "草绿色", "花色", "灰绿色", "奶白色", "银色", "黑灰色", "黄绿色", "粉白色", "嫩绿色", "奶金色", "浅绿色", "绛紫色", "浅紫色", "豆绿色", "黄琥珀色", "红色"}
var lpColorHair = []string{"橙色", "湖蓝色", "橘黄色", "深红色", "淡蓝色", "宝石蓝", "蓝白色", "深绿色", "浅绿色", "鹅黄色", "湖蓝色", "蓝黑色", "黄琥珀色", "浅紫色", "灰色", "灰绿色", "黄绿色", "蓝绿色", "海藻绿", "嫩绿色", "豆芽绿色", "豆绿色", "点缀着星空的夜之色", "浅橙色", "橙红色", "黑白闪烁", "蓝红渐变", "深蓝色", "蓝色", "可怜的原谅色", "草绿色", "棕红色", "西瓜红", "紫色", "浅紫色", "棕红色", "纯黑色", "黑灰色", "灰色", "白色", "蜜糖色", "枫叶色", "奶白色", "花色", "红色", "棕色", "墨绿色", "金色", "奶金色", "淡金色", "碎金色", "绛紫色", "粉白色", "稻草色", "紫黄渐变", "银色", "青山黛", "罗兰紫", "星河银", "翡冷翠", "夜光黑", "丹霞橙"}
var lpHairstyle = []string{"马尾", "高马尾", "侧单马尾", "半马尾", "双马尾", "双螺旋", "披肩双马尾", "四马尾", "多马尾", "麻花辫", "包子头", "朝天辫", "盘发", "辫子", "环形辫", "尾扎长发", "公主辫", "Half-up", "王冠编发", "翻翘", "猫耳型", "直发", "姬发式"}
var lpFeature = []string{"\n一对可爱的虎牙时不时露出来", "\n一条修长的龙尾彰显着不凡", "\n一对猫耳在头上时不时抖动一下", "\n一对兔耳跟着脑袋晃动着", "\n身上隐隐透出一股龙的威严！？", "\n一条尾巴在身后晃来晃去", "\n头顶有一个奇异的光环", "\n身后的小翅膀呼扇呼扇", "\n头上的双角宣示着小恶魔的身份", "\n头顶一根呆毛一晃一晃", "\n身后洁白的羽翼仿佛像天使", "\n尖尖的耳朵一动一动", "\n一颗泪痣点缀在眼角下方", "\n怎么看上去好像是隔壁神社的狐仙大人？？？", "\n一对小鹿角从头发中探出来", "\n头上有一对长长的天线", "\n一对蝙蝠翼从背后展开", "\n一条美丽的鱼尾吸引着你的目光", "\n长长的马尾随风飘动，如同舞者的轻盈", "\n身上覆盖着细密的鹰羽", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}
var lpOppai = []string{"一片平坦残念的样子", "看上去很平但是实际颇为有料", "一对硕大，仿佛吸住了你的眼睛", "毫无波澜的坦途", "穿着坚硬的板甲", "大小刚刚好", "似乎垫着什么东西", "一片平坦残念的样子", "看上去很平但是实际颇为有料", "一对硕大，仿佛吸住了你的眼睛", "毫无波澜的坦途", "大小刚刚好", "似乎垫着什么东西", "一片平坦残念的样子", "看上去很平但是实际颇为有料", "一对硕大，仿佛吸住了你的眼睛", "毫无波澜的坦途", "大小刚刚好", "似乎垫着什么东西"}
var lpSkin = []string{
    "柔嫩的奶白色",
    "有些不健康的苍白色",
    "有些害羞的粉白色",
    "健康的小麦色",
    "涉谷系的晒黑色",
    "透明般的白皙色",
    "温润的象牙色",
    "阳光下微微泛红的健康肤色",
    "冷艳的冰雪色",
    "瓷娃娃般的玉白色",
    "珍珠般的光泽肤色",
    "略带粉色的牛奶肤色",
}
var lpPersonality = []string{"傲娇", "傲沉", "病娇", "病切", "弱娇", "冷娇", "酷娇", "率直冷", "凛娇", "暴娇", "郁娇", "话痨", "三无", "无口", "冰美人", "高岭之花", "不悯", "KY", "无铁炮", "搞事", "暴力女", "高飞车", "笨蛋", "天然呆", "天然疯", "天然黑", "天然萌", "天真无邪", "电波", "冒失", "路痴", "大和抚子", "工口(糟糕)", "毒舌", "元气", "强气", "弱气", "怕羞", "忠犬", "爱哭鬼", "热血少女", "文学少女", "运动少女", "天才", "腹黑", "女王", "多重人格", "无节操", "口嫌体正直", "S属性", "M属性", "治愈系", "女神系", "残念系", "小恶魔系", "宝冢系", "不可思议系", "反差萌", "蠢萌", "节能主义", "拜金", "败家", "别扭", "单纯", "怀旧", "孤僻", "黑化", "叔萝"}
var lpRelation = []string{"恋人未满", "恋人", "未婚妻", "前任交往对象", "青梅竹马(幼驯染)", "初恋对象", "暗恋对象", "宿敌", "幻想女友", "学妹", "学姐", "同班同学", "同桌", "学生", "老师", "同事", "上司", "合租对象", "邻居", "双胞胎妹妹", "双胞胎姐姐", "义妹", "表妹", "堂妹", "义姐", "表姐", "堂姐", "女仆"}

var enType = map[string]string{
    "幼女":      "loli",
    "萝莉":      "loli",
    "万年萝莉":    "loli",
    "合法萝莉":    "loli",
    "萝莉老太婆": "loli",
    "乙女":      "bishoujo",
    "御姐":      "milf",
    "非法御姐":   "milf",
    "软妹":      "adolescent",
    "熟女":      "milf",
    "人妻":      "milf",
    "幼妻":      "loli",
    "BBA":     "BBA",
    "伪娘":      "crossdressing",
    "秀吉":      "crossdressing",
    "假小子":     "tomboy",
    "伪伪娘": "tomboy",
}

var enColor = map[string]string{
    "橙色":      "orange",
    "湖蓝色":     "lake blue",
    "橘黄色":     "orange-yellow",
    "深红色":     "deep red",
    "淡蓝色":     "light blue",
    "宝石蓝":     "jewel blue",
    "蓝白色":     "blue and white",
    "深绿色":     "dark green",
    "浅绿色":     "light green",
    "鹅黄色":     "light yellow",
    "蓝黑色":     "blue-black",
    "黄琥珀色":    "amber",
    "浅紫色":     "light purple",
    "灰色":      "gray",
    "灰绿色":     "gray-green",
    "黄绿色":     "yellow-green",
    "蓝绿色":     "blue-green",
    "海藻绿":     "Seagreen",
    "嫩绿色":     "tender green",
    "豆芽绿色":    "bean-green",
    "豆绿色":     "bean-green",
    "点缀着星空的夜之色": "black",
    "浅橙色":     "light orange",
    "橙红色":     "orange-red",
    "黑白闪烁":    "black and white",
    "蓝红渐变":    "blue-red gradient",
    "深蓝色":     "dark blue",
    "蓝色":      "blue",
    "可怜的原谅色":  "green",
    "草绿色":     "grass green",
    "棕红色":     "brown-red",
    "西瓜红":     "red",
    "紫色":      "purple",
    "纯黑色":     "pure black",
    "黑灰色":     "black gray",
    "白色":      "white",
    "蜜糖色":     "golden warmth",
    "枫叶色":     "warmth red",
    "奶白色":     "milky-white",
    "花色":      "floral",
    "红色":      "red",
    "棕色":      "brown",
    "墨绿色":     "dark green",
    "金色":      "gold",
    "奶金色":     "light gold",
    "淡金色":     "light gold",
    "碎金色":     "broken gold",
    "绛紫色":     "red-purple",
    "粉白色":     "pink white",
    "稻草色":     "straw color",
    "紫黄渐变":    "purple-yellow gradient",
    "银色":      "silver",
    "青山黛":     "cyan",
    "罗兰紫":     "purple",
    "星河银":     "silver",
    "翡冷翠":     "green",
    "夜光黑":     "black",
    "丹霞橙":     "orange",
}

var enHairstyle = map[string]string{
    "马尾":      "ponytail",
    "高马尾":     "high ponytail",
    "侧单马尾":    "side ponytail",
    "半马尾":     "half ponytail",
    "双马尾":     "twintails",
    "双螺旋":     "double helix",
    "披肩双马尾":   "medium hair, twintails",
    "四马尾":     "quad ponytails hair",
    "多马尾":     "multiple ponytails hair",
    "麻花辫":     "braid",
    "包子头":     "hair bun",
    "朝天辫":     "upward-braid",
    "盘发":      "hair up",
    "辫子":      "braid",
    "环形辫":     "circular braid",
    "尾扎长发":    "tie long hair",
    "公主辫":     "princess braid",
    "Half-up":  "Half-up",
    "王冠编发":    "crown braid",
    "翻翘":      "flipped",
    "猫耳型": "cat ear style hair",
    "直发": "straight hair",
    "姬发式": "hime cut",
}

var enFeature = map[string]string{
    "\n一对可爱的虎牙时不时露出来": "skin fang",
    "\n头顶一根呆毛一晃一晃": "ahoge",
    "\n一对猫耳在头上时不时抖动一下": "cat ears, cat girl, cat tail",
    "\n一对兔耳跟着脑袋晃动着": "rabbit ears, rabbit girl, rabbit tail",
    "\n一颗泪痣点缀在眼角下方": "mole under eye",
    "\n一条尾巴在身后晃来晃去": "dog ears, dog tail, dog girl, dog collar",
    "\n头顶有一个奇异的光环": "halo",
    "\n身后的小翅膀呼扇呼扇": "small wings",
    "\n头上的双角宣示着小恶魔的身份": "succubus, demon horn, demon tail",
    "\n身后洁白的羽翼仿佛像天使": "angel, angel wings, halo",
    "\n尖尖的耳朵一动一动": "pointy ears, elf",
    "\n身上隐隐透出一股龙的威严！？": "dragon girl, dragon wings, dragon horns, dragon tail",
    "\n怎么看上去好像是隔壁神社的狐仙大人？？？": "fox girl, fox tail, fox ears",
    "\n一对小鹿角从头发中探出来": "deer ears, antlers", 
    "\n头上有一对长长的天线": "gynoid, mecha musume, antenna", 
    "\n一对蝙蝠翼从背后展开": "vampire, bat wings, skin fangs", 
    "\n一条美丽的鱼尾吸引着你的目光": "mermaid, head fins, fish tail, fish hair ornament",
    "\n长长的马尾随风飘动，如同舞者的轻盈": "pony girl, horse ears, horse tail",
    "\n身上覆盖着细密的鹰羽": "(harpy, harpie lady, winged arms, claws)",
    "\n一条修长的龙尾彰显着不凡": "(eastern dragon girl, eastern dragon horns, eastern dragon tail)",
}

var enOppai = map[string]string{
    "一片平坦残念的样子": "flat chest",
    "看上去很平但是实际颇为有料": "small breasts",
    "一对硕大，仿佛吸住了你的眼睛": "big breasts",
    "毫无波澜的坦途": "flat chest",
    "穿着坚硬的板甲": "armor, shoulder armor, armored dress",
    "大小刚刚好": "medium breasts",
    "似乎垫着什么东西": "medium breasts",
}

var enSkin = map[string]string{
    "有些不健康的苍白色": "pale skin",
    "健康的小麦色": "tan",
    "涉谷系的晒黑色": "tan",
    "透明般的白皙色": "white skin",
    "阳光下微微泛红的健康肤色": "tan",
    "冷艳的冰雪色": "pale skin",
    "瓷娃娃般的玉白色": "shiny skin",
    "珍珠般的光泽肤色": "shiny skin",
}

var enPersonality = map[string]string{
    "傲娇": "tsundere",
    "傲沉": "hinedere",
    "病娇": "yandere",
    "病切": "byoukidere",
    "弱娇": "byoukidere",
    "冷娇": "kuudere",
    "酷娇": "kuudere",
    "率直冷": "cool",
    "凛娇": "tsundere",
    "暴娇": "Bokodere",
    "郁娇": "darudere",
    "三无": "dandere",
    "无口": "dandere",
    "冰美人": "Ice Beauty",
    "高岭之花": "Ice Beauty",
    "无铁炮": "clumsy",
    "搞事": "violent girl",
    "暴力女": "violent girl",
    "高飞车": "violent girl",
    "笨蛋": "bakadere",
    "天然呆": "innocent",
    "天然疯": "innocent",
    "天然黑": "innocent",
    "天然萌": "innocent",
    "天真无邪": "innocent",
    "冒失": "clumsy",
    "大和抚子": "Yamato Nadeshiko",
    "工口(糟糕)": "female pervert",
    "元气": "energetic",
    "强气": "strong-willed",
    "弱气": "weak-willed",
    "怕羞": "shy",
    "爱哭鬼": "Crybaby",
    "文学少女": "bungaku shoujo",
    "运动少女": "sports wear",
    "女王": "queen",
    "无节操": "female pervert",
    "S属性": "Sadodere, S dere",
    "M属性": "masodere, M dere",
    "小恶魔系": "imp",
    "宝冢系": "gyaru",
    "别扭": "embarrased",
    "单纯": "simple",
    "孤僻": "lonely",
    "黑化": "evil smile",
    "叔萝": "loli",
}

var enRelation = map[string]string{
    "学妹": "student", "学姐": "student", "同班同学": "student", 
    "同桌": "student", "学生": "student", "老师": "teacher", "女仆": "maid", "巫女": "miko",
}


// 获取随机元素
func getRandomElement(ctx *zero.Ctx,list []string) string {
    seedStr := strconv.FormatInt(ctx.Event.UserID, 10) + time.Now().Format("2006-01-02")
    h := fnv.New64a()
    h.Write([]byte(seedStr))
    seed := int64(h.Sum64())
    rand.Seed(seed)
    return list[rand.Intn(len(list))]
}

// 生成角色描述
func generateWaifu(ctx *zero.Ctx) (string, string) {
    lpTypeVal := getRandomElement(ctx,lpType)
    lpColorHairVal := getRandomElement(ctx,lpColorHair)
    lpColorEyesVal := getRandomElement(ctx,lpColorEyes)
    lpHairstyleVal := getRandomElement(ctx,lpHairstyle)
    lpFeatureVal := getRandomElement(ctx,lpFeature)
    lpOppaiVal := getRandomElement(ctx,lpOppai)
    lpSkinVal := getRandomElement(ctx,lpSkin)
    lpPersonalityVal := getRandomElement(ctx,lpPersonality)
    lpRelationVal := getRandomElement(ctx,lpRelation)

    chineseDescription := fmt.Sprintf("你感觉到空气中一阵凝滞，似乎整个周遭仿佛成为了一张巨画，铺展开来。\n"+
        "你遭遇了降维打击，等你醒过来，你发现一位少女正站在你的面前，你知道，这是你的纸片人老婆。\n"+
        "她是一个%s\n有着一头%s的%s\n眼睛是%s的%s\n胸前%s\n肤色是%s的\n看起来挺%s的样子\n实际上是你的%s",
        lpTypeVal, lpColorHairVal, lpHairstyleVal, lpColorEyesVal, lpFeatureVal, lpOppaiVal, lpSkinVal, lpPersonalityVal, lpRelationVal)

    englishDescription := fmt.Sprintf("%s, %s hair, %s, %s eyes, %s, %s, %s, %s, %s",
        enType[lpTypeVal], enColor[lpColorHairVal], enHairstyle[lpHairstyleVal], enColor[lpColorEyesVal], enFeature[lpFeatureVal], 
        enOppai[lpOppaiVal], enSkin[lpSkinVal], enPersonality[lpPersonalityVal], enRelation[lpRelationVal])

    return chineseDescription, englishDescription
}

func genprompt(ctx *zero.Ctx) string {
    chineseDesc, englishDesc := generateWaifu(ctx)
    ctx.SendChain(
        message.At(ctx.Event.UserID),
        message.Text(chineseDesc),
    )
	return englishDesc
}


// RequestData 定义发送给API的数据结构
type RequestData struct {
    Prompt       string   `json:"prompt"`
    Steps        int      `json:"steps"`
    Width        int      `json:"width"`
    Height       int      `json:"height"`
    Sampler      string   `json:"sampler_name"`
    Styles       []string `json:"styles"`
    Enable_hr    string   `json:"enable_hr"`
    CFG_scale    string   `json:"cfg_scale"`
    Denoising_strength string `json:"denoising_strength"`
    Clip_skip    int   `json:"clip_skip"`
}

func init() { //插件主体

	engine := control.AutoRegister(&ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "纸片人老婆",
		Help:             "纸片人老婆",
		PublicDataFolder: "2dwife",
	}).ApplySingle(ctxext.DefaultSingle)

	engine.OnFullMatchGroup([]string{"纸片人老婆"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			    // 接受prompt
    prompt := genprompt(ctx)
    prompt = "1girl, solo, " + prompt

    // 创建请求数据
    requestData := RequestData{
        Prompt:  prompt,
        Steps:   20,  // 步数
        Width:   360, // 图片宽度
        Height:  540, // 图片高度
        Sampler: "Euler a", //采样器
        Styles:  []string{"cqhttp"}, //styles
        Enable_hr: "true", //高清修复
        Denoising_strength: "0.7", //修复幅度
        CFG_scale: "7", //CFG
        Clip_skip: 2,
    }

    // 将请求数据编码为JSON
    jsonData, err := json.Marshal(requestData)
    if err != nil {
        fmt.Println("JSON编码失败:", err)
        return
    }

    // 定义API端点
    url := "http://127.0.0.1:7860/sdapi/v1/txt2img" //本地sdwebui的api地址

    // 创建HTTP请求
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("创建HTTP请求失败:", err)
        return
    }
    req.Header.Set("Content-Type", "application/json")

    // 发送HTTP请求
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("发送HTTP请求失败:", err)
        return
    }
    defer resp.Body.Close()

    // 读取响应
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("读取响应失败:", err)
        return
    }

    // 打印响应
    fmt.Println("API响应:", string(body))

    // 解析响应中的图像数据
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        fmt.Println("JSON解析失败:", err)
        return
    }

    images, ok := result["images"].([]interface{})
    if !ok || len(images) == 0 {
        fmt.Println("没有生成图像")
        return
    }

	// 获取图片并发送给用户
	imgBase64 := images[0].(string)
	imgData, err := base64.StdEncoding.DecodeString(imgBase64)
	if err != nil {
		fmt.Println("图像数据解码失败:", err)
		return
	}

	// 发送图片给用户
	ctx.SendChain(
		message.At(ctx.Event.UserID),
		message.ImageBytes(imgData),
	)
		})
}
