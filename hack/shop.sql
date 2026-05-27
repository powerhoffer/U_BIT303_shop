/*
 Navicat Premium Data Transfer

 Source Server         : 本地
 Source Server Type    : MySQL
 Source Server Version : 80012
 Source Host           : localhost:3306
 Source Schema         : shop

 Target Server Type    : MySQL
 Target Server Version : 80012
 File Encoding         : 65001

 Date: 05/12/2025 16:07:29
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for address_info
-- ----------------------------
DROP TABLE IF EXISTS `address_info`;
CREATE TABLE `address_info`  (
  `id` int(11) NOT NULL,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `pid` int(11) NOT NULL,
  `status` int(11) NOT NULL DEFAULT 0,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `pid`(`pid` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '全国城市信息表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of address_info 只保留北京的作为演示
-- ----------------------------
INSERT INTO `address_info` VALUES (110000, '北京', 1, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110100, '北京市', 110000, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110101, '东城区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110102, '西城区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110103, '崇文区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110104, '宣武区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110105, '朝阳区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110106, '丰台区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110107, '石景山区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110108, '海淀区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110109, '门头沟区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110111, '房山区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110112, '通州区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110113, '顺义区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110114, '昌平区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110115, '大兴区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110116, '怀柔区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110117, '平谷区', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110228, '密云县', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110229, '延庆县', 110100, 0, '2013-07-10 11:02:58');
INSERT INTO `address_info` VALUES (110230, '其它区', 110100, 0, '2013-07-10 11:02:58');

-- ----------------------------
-- Table structure for admin_info
-- ----------------------------
DROP TABLE IF EXISTS `admin_info`;
CREATE TABLE `admin_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `role_ids` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '角色ids',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `user_salt` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '加密盐',
  `is_admin` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否超级管理员',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name_unique`(`name` ASC) USING BTREE COMMENT '名字唯一索引'
) ENGINE = InnoDB AUTO_INCREMENT = 26 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of admin_info
-- ----------------------------
INSERT INTO `admin_info` VALUES (1, 'wangzhongyang', 'ce1d477c2c7a11384053d8981d85bd96', '2', '2025-12-04 08:21:39', '2025-12-04 09:00:18', '9IQDuhRTJW', 1);
INSERT INTO `admin_info` VALUES (2, 'test', '9122185be9f038b36b30bd65daa79b80', '1', '2025-12-04 08:42:32', '2025-12-04 08:42:32', 'vId0yaMHzQ', 0);

-- ----------------------------
-- Table structure for article_info
-- ----------------------------
DROP TABLE IF EXISTS `article_info`;
CREATE TABLE `article_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT 0 COMMENT '作者id',
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标题',
  `desc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '摘要',
  `pic_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '封面图',
  `is_admin` tinyint(1) NOT NULL DEFAULT 2 COMMENT '1后台管理员发布 2前台用户发布',
  `praise` int(11) NOT NULL DEFAULT 0 COMMENT '点赞数',
  `detail` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '文章详情',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '文章（种草）表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of article_info
-- ----------------------------
<<<<<<< HEAD
-- 华凌空调种草文章 - 真实用户体验
INSERT INTO `article_info` VALUES (1, 101, '华凌N8HE1空调使用3个月真实感受：卧室制冷神器', '租房党福音！静音省电，制冷快，这个夏天终于能睡好觉了', 'https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fgfs17.gomein.net.cn%2FT108VWB4W_1RCvBVdK_800.jpg%3Fv%3D1&refer=http%3A%2F%2Fgfs17.gomein.net.cn&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1660794257&t=795ee536d5af33788a249b08d0b28b6f', 2, 236, '坐标上海，卧室15平左右，6月份入手的华凌N8HE1。安装师傅很专业，半小时搞定。

使用感受：
1. 制冷真的快！开机5分钟整个房间就凉下来了，温度很均匀
2. 超级静音！晚上开睡眠模式，几乎听不到声音，再也不会被空调噪音吵醒
3. 省电！一个月下来电费比去年用老空调省了快100块
4. 外观简约大气，白色百搭，和我家装修风格很配

缺点：遥控器有点简单，要是能连接手机APP就更好了。不过这个价格，还要什么自行车！', '2022-06-15 14:30:22', '2022-06-15 14:30:22', NULL);

-- 小米空气净化器种草文章 - 过敏患者福音
INSERT INTO `article_info` VALUES (2, 102, '小米空气净化器Pro H：过敏体质的救星', '养宠物+过敏体质的我，终于找到了呼吸自由的秘诀', 'https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/5f7a9c5b9e1c4a1d0d3d3d3d3d3d3d3d.jpg', 2, 189, '家里养了两只猫，我又有过敏性鼻炎，一到春秋就难受得不行。朋友推荐了小米空气净化器Pro H，用了半个月，效果真的惊艳！

使用体验：
1. 净化速度快！开机半小时，PM2.5从120降到20以下
2. 智能感应很灵敏，宠物掉毛多的时候会自动调高风速
3. 静音效果好，晚上开睡眠模式几乎没声音
4. 手机APP可以实时查看空气质量，还能远程控制

现在每天起床鼻子都很舒服，再也不会一直打喷嚏了。强烈推荐给养宠物和过敏体质的朋友！', '2022-09-28 09:15:47', '2022-09-28 09:15:47', NULL);

-- 戴森吸尘器种草文章 - 家居清洁神器
INSERT INTO `article_info` VALUES (3, 103, '戴森V12吸尘器：让做家务变成一种享受', '无线轻便，吸力强，各种缝隙都能轻松搞定', 'https://www.dyson.cn/content/dam/dyson-cn/images/products/cordless-vacuums/v12/dyson-v12-detect-slim/fluffy-main-image.png', 2, 452, '作为一个爱干净的家庭主妇，试过很多吸尘器，直到遇到戴森V12。

使用感受：
1. 无线设计太方便了！不用再拖着线到处跑，想吸哪里吸哪里
2. 吸力真的强！地板上的头发、灰尘，甚至地毯里的深层污垢都能吸干净
3. 配件丰富，沙发、窗帘、床褥都能用不同的吸头搞定
4. 续航不错，120平的房子一次吸完还有电
5. 自清洁功能很贴心，不用手动清理滤网

虽然价格有点贵，但真的物超所值，大大节省了做家务的时间！', '2022-11-12 16:45:12', '2022-11-12 16:45:12', NULL);

-- 美的冰箱种草文章 - 三口之家的选择
INSERT INTO `article_info` VALUES (4, 104, '美的BCD-478WSPZM(E)冰箱：大容量还能除味', '三口之家够用，保鲜效果好，再也没有异味困扰', 'https://img12.360buyimg.com/n7/jfs/t1/198818/30/30076/123456/632c9a29E1a7a5a5a/1a1a1a1a1a1a1a1a.jpg', 2, 156, '家里之前的冰箱太小了，放不了多少东西，还容易串味。双十一换了美的这款478升的冰箱，真的太香了！

使用体验：
1. 容量大！冷藏冷冻都很能装，三口之家完全够用
2. 保鲜效果好！放了一周的蔬菜还是很新鲜，水果也不容易坏
3. 除味效果不错，打开冰箱再也没有难闻的异味
4. 外观是磨砂玻璃，很高级，容易清洁
5. 一级能效，省电又静音

价格也很亲民，性价比很高，推荐给需要换冰箱的家庭！', '2022-12-05 11:20:33', '2022-12-05 11:20:33', NULL);

-- 索尼电视种草文章 - 影音爱好者必入
INSERT INTO `article_info` VALUES (5, 105, '索尼XR-65X90J电视：在家就能享受影院级体验', '画质惊艳，音效震撼，看电影玩游戏都超爽', 'https://www.sony.com.cn/image/0014b557b6d34b99a52b39a7000d3d3d?fmt=pjpeg&wid=1200&hei=630&bgcolor=FFFFFF&bgc=FFFFFF', 2, 318, '作为一个影音爱好者，纠结了很久终于入手了索尼XR-65X90J。

使用感受：
1. 画质真的绝了！色彩鲜艳真实，对比度高，HDR效果惊艳
2. 音效震撼，自带的扬声器就有很好的环绕感
3. 系统流畅，反应快，没有广告
4. 游戏模式延迟很低，玩PS5体验超棒
5. 外观简约大气，挂墙效果很好

虽然价格不便宜，但对于追求画质和音效的朋友来说，绝对值得入手！', '2023-02-18 19:45:56', '2023-02-18 19:45:56', NULL);

-- 小熊加湿器种草文章 - 干燥季节必备
INSERT INTO `article_info` VALUES (6, 106, '小熊加湿器JSQ-C40L5：办公室和卧室都能用', '静音加湿，容量大，再也不怕皮肤干燥了', 'https://img14.360buyimg.com/n7/jfs/t1/187655/33/29876/123456/632c9a29E1a7a5a5a/2b2b2b2b2b2b2b2b.jpg', 2, 98, '北方的冬天真的太干燥了，办公室和家里都需要加湿器。对比了很多款，选了小熊这款。

使用体验：
1. 容量大！4升的水箱，加满水可以用一整晚
2. 加湿效果好，房间湿度能保持在40-50%之间
3. 静音！办公室用不会影响同事，晚上用也不会吵
4. 有定时功能，很方便
5. 外观可爱，放在桌面不占地方

价格便宜，质量不错，推荐给需要加湿器的朋友！', '2023-01-10 10:15:22', '2023-01-10 10:15:22', NULL);

-- 苏泊尔电饭煲种草文章 - 米饭爱好者的福音
INSERT INTO `article_info` VALUES (7, 107, '苏泊尔SF40HC88电饭煲：煮出的米饭真的不一样', 'IH加热，口感Q弹，各种米都能煮出好味道', 'https://img13.360buyimg.com/n7/jfs/t1/198818/30/30076/123456/632c9a29E1a7a5a5a/3c3c3c3c3c3c3c3c.jpg', 2, 175, '家里之前的电饭煲用了5年，煮出来的米饭总是软硬不均。换了苏泊尔这款IH电饭煲，真的惊艳到我了！

使用体验：
1. 米饭口感真的好！Q弹有嚼劲，不管是东北米还是泰国香米都能煮出好味道
2. 功能丰富，有柴火饭、粥、汤、蛋糕等多种模式
3. 预约功能很方便，早上出门前预约，晚上回家就能吃现成的
4. 内胆厚重，导热均匀，不容易粘锅
5. 外观简约大气，容易清洁

价格适中，性价比很高，推荐给注重米饭口感的朋友！', '2023-03-25 15:30:44', '2023-03-25 15:30:44', NULL);

-- 苹果AirPods Pro 2种草文章 - 通勤必备神器
INSERT INTO `article_info` VALUES (8, 108, 'AirPods Pro 2使用体验：降噪效果真的能打', '通勤路上的救星，音质好，佩戴舒适', 'https://store.storeimages.cdn-apple.com/4982/as-images.apple.com/is/MQD83CH/A_AV1?wid=2000&hei=2000&fmt=jpeg&qlt=95&.v=1664472259559', 2, 422, '作为一个每天通勤2小时的打工人，AirPods Pro 2真的是我的救星！

使用体验：
1. 降噪效果太绝了！地铁上的噪音几乎完全听不到，世界瞬间安静
2. 音质比第一代有提升，低音更饱满，高音清晰
3. 佩戴舒适，长时间戴也不会耳朵疼
4. 续航不错，单次能用6小时，配合充电盒能用30小时
5. 空间音频效果惊艳，看电影就像在电影院一样

虽然价格不便宜，但对于每天通勤的人来说，真的能提升生活质量！', '2023-04-12 18:55:17', '2023-04-12 18:55:17', NULL);
-- ----------------------------
-- Table structure for cart_info
-- ----------------------------
DROP TABLE IF EXISTS `cart_info`;
CREATE TABLE `cart_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '购物车表',
  `user_id` int(11) NOT NULL DEFAULT 0,
  `goods_options_id` int(11) NOT NULL DEFAULT 0 COMMENT '商品规格id',
  `count` int(11) NOT NULL COMMENT '商品数量',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of cart_info
-- ----------------------------
INSERT INTO `cart_info` VALUES (1, 1, 2, 1, '2022-07-29 13:59:10', '2022-07-29 13:59:10', NULL);
INSERT INTO `cart_info` VALUES (2, 1, 8, 1, '2022-07-29 14:23:31', '2022-07-29 14:23:31', NULL);
INSERT INTO `cart_info` VALUES (3, 1, 11, 2, '2022-07-29 14:30:00', '2022-07-29 14:30:00', NULL);

-- ----------------------------
-- Table structure for category_info
-- ----------------------------
DROP TABLE IF EXISTS `category_info`;
CREATE TABLE `category_info`  (
  `id` int(0) NOT NULL AUTO_INCREMENT,
  `parent_id` int(0) NOT NULL DEFAULT 0 COMMENT '父级id',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分类名称',
  `pic_url` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分类图标',
  `level` tinyint(1) NOT NULL DEFAULT 1 COMMENT '等级 默认1级分类',
  `sort` int(0) NOT NULL DEFAULT 1 COMMENT '排序值',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  `description` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分类描述',
  `created_at` datetime(0) DEFAULT NULL,
  `updated_at` datetime(0) DEFAULT NULL,
  `deleted_at` datetime(0) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_parent_id`(`parent_id`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '商品分类表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of category_info
-- ----------------------------
-- 一级分类
INSERT INTO `category_info` VALUES (1, 0, '家用电器', 'https://img.yzcdn.cn/vant/apple-1.jpg', 1, 1, 1, '涵盖各类家电产品，满足家庭生活需求', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (2, 0, '手机数码', 'https://img.yzcdn.cn/vant/apple-2.jpg', 1, 2, 1, '手机、数码配件等智能产品', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (3, 0, '电脑办公', 'https://img.yzcdn.cn/vant/apple-3.jpg', 1, 3, 1, '电脑、办公设备及耗材', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (4, 0, '家居家装', 'https://img.yzcdn.cn/vant/apple-4.jpg', 1, 4, 1, '家居用品、装修材料等', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (5, 0, '服饰鞋包', 'https://img.yzcdn.cn/vant/apple-5.jpg', 1, 5, 1, '服装、鞋靴、箱包等时尚单品', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- 二级分类 - 家用电器
INSERT INTO `category_info` VALUES (6, 1, '电视', 'https://img.yzcdn.cn/vant/apple-1.jpg', 2, 1, 1, '各类智能电视产品', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (7, 1, '空调', 'https://img.yzcdn.cn/vant/apple-1.jpg', 2, 2, 1, '空调制冷制热设备', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (8, 1, '冰箱', 'https://img.yzcdn.cn/vant/apple-1.jpg', 2, 3, 1, '冷藏冷冻保鲜设备', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (9, 1, '洗衣机', 'https://img.yzcdn.cn/vant/apple-1.jpg', 2, 4, 1, '衣物清洗设备', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (10, 1, '厨房电器', 'https://img.yzcdn.cn/vant/apple-1.jpg', 2, 5, 1, '厨房用各类电器', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- 二级分类 - 手机数码
INSERT INTO `category_info` VALUES (11, 2, '手机', 'https://img.yzcdn.cn/vant/apple-2.jpg', 2, 1, 1, '智能手机产品', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (12, 2, '耳机', 'https://img.yzcdn.cn/vant/apple-2.jpg', 2, 2, 1, '各类耳机产品', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (13, 2, '智能穿戴', 'https://img.yzcdn.cn/vant/apple-2.jpg', 2, 3, 1, '智能手表、手环等', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (14, 2, '摄影摄像', 'https://img.yzcdn.cn/vant/apple-2.jpg', 2, 4, 1, '相机、摄像机等设备', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- 二级分类 - 电脑办公
INSERT INTO `category_info` VALUES (15, 3, '笔记本电脑', 'https://img.yzcdn.cn/vant/apple-3.jpg', 2, 1, 1, '便携式笔记本电脑', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (16, 3, '台式机', 'https://img.yzcdn.cn/vant/apple-3.jpg', 2, 2, 1, '台式电脑主机', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (17, 3, '平板电脑', 'https://img.yzcdn.cn/vant/apple-3.jpg', 2, 3, 1, '平板设备', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (18, 3, '办公耗材', 'https://img.yzcdn.cn/vant/apple-3.jpg', 2, 4, 1, '打印纸、墨盒等耗材', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- 三级分类 - 电视
INSERT INTO `category_info` VALUES (19, 6, '全面屏电视', 'https://img.yzcdn.cn/vant/apple-1.jpg', 3, 1, 1, '无边框全面屏设计电视', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (20, 6, '4K超高清电视', 'https://img.yzcdn.cn/vant/apple-1.jpg', 3, 2, 1, '4K分辨率超高清电视', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (21, 6, '8K超高清电视', 'https://img.yzcdn.cn/vant/apple-1.jpg', 3, 3, 1, '8K分辨率超高清电视', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (22, 6, '量子点电视', 'https://img.yzcdn.cn/vant/apple-1.jpg', 3, 4, 1, '量子点显示技术电视', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (23, 6, 'OLED电视', 'https://img.yzcdn.cn/vant/apple-1.jpg', 3, 5, 1, 'OLED自发光显示电视', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- 三级分类 - 空调
INSERT INTO `category_info` VALUES (24, 7, '挂壁式空调', 'https://img.yzcdn.cn/vant/apple-1.jpg', 3, 1, 1, '卧室客厅挂壁安装空调', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (25, 7, '立柜式空调', 'https://img.yzcdn.cn/vant/apple-1.jpg', 3, 2, 1, '大空间立柜式空调', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (26, 7, '中央空调', 'https://img.yzcdn.cn/vant/apple-1.jpg', 3, 3, 1, '多房间中央空调系统', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (27, 7, '风管机', 'https://img.yzcdn.cn/vant/apple-1.jpg', 3, 4, 1, '隐藏式风管空调', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- 三级分类 - 手机
INSERT INTO `category_info` VALUES (28, 11, '5G手机', 'https://img.yzcdn.cn/vant/apple-2.jpg', 3, 1, 1, '支持5G网络的智能手机', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (29, 11, '折叠屏手机', 'https://img.yzcdn.cn/vant/apple-2.jpg', 3, 2, 1, '可折叠屏幕智能手机', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (30, 11, '游戏手机', 'https://img.yzcdn.cn/vant/apple-2.jpg', 3, 3, 1, '高性能游戏专用手机', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (31, 11, '拍照手机', 'https://img.yzcdn.cn/vant/apple-2.jpg', 3, 4, 1, '高像素拍照智能手机', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- 三级分类 - 笔记本电脑
INSERT INTO `category_info` VALUES (32, 15, '轻薄本', 'https://img.yzcdn.cn/vant/apple-3.jpg', 3, 1, 1, '轻薄便携笔记本电脑', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (33, 15, '游戏本', 'https://img.yzcdn.cn/vant/apple-3.jpg', 3, 2, 1, '高性能游戏笔记本电脑', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (34, 15, '商务本', 'https://img.yzcdn.cn/vant/apple-3.jpg', 3, 3, 1, '商务办公专用笔记本', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `category_info` VALUES (35, 15, '二合一笔记本', 'https://img.yzcdn.cn/vant/apple-3.jpg', 3, 4, 1, '平板笔记本二合一设备', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- ----------------------------
-- Table structure for collection_info
-- ----------------------------
DROP TABLE IF EXISTS `collection_info`;
CREATE TABLE `collection_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT 0 COMMENT '用户id',
  `object_id` int(11) NOT NULL DEFAULT 0 COMMENT '对象id',
  `type` tinyint(1) NOT NULL DEFAULT 0 COMMENT '收藏类型：1商品 2文章',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_index`(`user_id` ASC, `object_id` ASC, `type` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of collection_info
-- ----------------------------
INSERT INTO `collection_info` VALUES (1, 1, 1, 1, '2022-07-31 15:21:38', '2022-07-31 15:21:38');
INSERT INTO `collection_info` VALUES (2, 1, 3, 1, '2022-07-31 15:22:00', '2022-07-31 15:22:00');
INSERT INTO `collection_info` VALUES (3, 1, 4, 1, '2022-07-31 15:23:00', '2022-07-31 15:23:00');
INSERT INTO `collection_info` VALUES (4, 1, 1, 2, '2022-07-31 15:24:00', '2022-07-31 15:24:00');
INSERT INTO `collection_info` VALUES (5, 1, 2, 2, '2022-07-31 15:25:00', '2022-07-31 15:25:00');

-- ----------------------------
-- Table structure for comment_info
-- ----------------------------
DROP TABLE IF EXISTS `comment_info`;
CREATE TABLE `comment_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) NOT NULL DEFAULT 0 COMMENT '父级评论id',
  `user_id` int(11) NOT NULL DEFAULT 0,
  `object_id` int(11) NOT NULL DEFAULT 0,
  `type` tinyint(1) NOT NULL DEFAULT 0 COMMENT '评论类型：1商品 2文章',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '评论内容',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_index`(`user_id`, `object_id`, `type`, `content`, `parent_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;


-- ----------------------------
-- Records of comment_info
-- ----------------------------
INSERT INTO `comment_info` VALUES (4, 0, 1, 1, 1, '商品质量非常好，与描述的完全一致，包装也很仔细，没有破损，物流速度很快，非常满意的一次购物！', '2022-07-31 17:23:48', '2022-07-31 17:23:48', NULL);
INSERT INTO `comment_info` VALUES (5, 0, 1, 2, 2, '这篇文章写得很详细，学到了很多知识，作者辛苦了，期待更多优质内容！', '2022-07-31 17:24:10', '2022-07-31 17:24:10', NULL);
INSERT INTO `comment_info` VALUES (6, 0, 2, 1, 1, '东西收到了，质量还可以，就是物流有点慢，等了好几天才到，总体来说还是不错的。', '2022-08-01 09:15:22', '2022-08-01 09:15:22', NULL);
INSERT INTO `comment_info` VALUES (7, 5, 1, 2, 2, '确实写得很好，对我帮助很大，感谢分享！', '2022-07-31 17:24:59', '2022-07-31 17:24:59', NULL);
INSERT INTO `comment_info` VALUES (8, 0, 3, 1, 1, '商品收到了，打开一看有点失望，和图片上的差距有点大，质量也一般般，不太推荐购买。', '2022-08-02 14:30:45', '2022-08-02 14:30:45', NULL);
INSERT INTO `comment_info` VALUES (9, 0, 4, 1, 1, '非常满意的一次购物，商品质量好，价格实惠，物流速度快，客服态度也很好，下次还会再来！', '2023-01-18 10:15:33', '2023-01-18 10:15:33', NULL);
INSERT INTO `comment_info` VALUES (10, 1, 4, 1, 1, '第二次购买了，质量还是一如既往的好，值得信赖！', '2023-01-19 14:25:24', '2023-01-19 14:25:24', NULL);
INSERT INTO `comment_info` VALUES (11, 3, 4, 1, 1, '物流速度很快，包装很严实，商品没有损坏，使用了几天感觉不错，性价比很高。', '2023-01-19 14:26:50', '2023-01-19 14:26:50', NULL);
INSERT INTO `comment_info` VALUES (12, 0, 5, 2, 2, '文章内容很有深度，分析得很到位，让我对这个问题有了新的认识，谢谢作者！', '2023-01-20 16:45:12', '2023-01-20 16:45:12', NULL);
INSERT INTO `comment_info` VALUES (13, 2, 5, 2, 2, '写得不错，但是有些地方可以再详细一点，期待后续更新。', '2023-01-21 09:30:28', '2023-01-21 09:30:28', NULL);
INSERT INTO `comment_info` VALUES (14, 0, 2, 2, 2, '内容很实用，按照文章的方法操作，确实解决了我的问题，非常感谢！', '2023-01-22 11:20:15', '2023-01-22 11:20:15', NULL);
INSERT INTO `comment_info` VALUES (15, 4, 3, 1, 1, '商品质量不错，就是价格有点贵，不过一分钱一分货，总体来说还是值得购买的。', '2023-01-23 15:45:36', '2023-01-23 15:45:36', NULL);

-- ----------------------------
-- Table structure for consignee_info
-- ----------------------------
DROP TABLE IF EXISTS `consignee_info`;
CREATE TABLE `consignee_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '收货地址表',
  `user_id` int(11) NOT NULL DEFAULT 0,
  `is_default` tinyint(1) NOT NULL DEFAULT 0 COMMENT '默认地址1  非默认0\n',
  `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `phone` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `province` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `city` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `town` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '县区',
  `street` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '街道乡镇',
  `detail` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '地址详情',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of consignee_info
-- ----------------------------
INSERT INTO `consignee_info` VALUES (1, 1, 1, '王中阳', '13800138000', '北京', '北京市', '朝阳区', '望京街道', '望京SOHO T1 1801室', '2022-07-31 14:42:33', '2022-07-31 14:44:50', NULL);
INSERT INTO `consignee_info` VALUES (2, 1, 0, '王中阳', '13800138000', '上海', '上海市', '浦东新区', '张江高科技园区', '博云路2号', '2022-07-31 14:45:00', '2022-07-31 14:45:00', NULL);

-- ----------------------------
-- Table structure for coupon_info
-- ----------------------------
DROP TABLE IF EXISTS `coupon_info`;
CREATE TABLE `coupon_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `price` int(11) NOT NULL DEFAULT 0 COMMENT '优惠前面值 单位分\n',
  `goods_ids` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '关联使用的goods_ids  逗号分隔',
  `category_id` int(11) NOT NULL DEFAULT 0 COMMENT '关联使用的分类id',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '轮播图表\n' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of coupon_info
-- ----------------------------
INSERT INTO `coupon_info` VALUES (1, '满2千减5百优惠券', 50000, '1,2,3', 1, '2022-07-19 14:30:48', '2022-09-25 15:35:56', NULL);
INSERT INTO `coupon_info` VALUES (2, '满2千减5百优惠券', 50000, '0', 1, '2022-07-19 14:39:51', '2022-07-19 14:39:51', NULL);
INSERT INTO `coupon_info` VALUES (3, '满2千减5百优惠券', 50000, '1', 1, '2022-07-29 15:58:15', '2022-08-01 13:53:11', '2022-08-01 13:53:27');
INSERT INTO `coupon_info` VALUES (4, '满2千减5百优惠券', 50000, '0', 1, '2022-08-01 13:52:51', '2022-08-01 13:52:51', NULL);
INSERT INTO `coupon_info` VALUES (5, '满2千减5百优惠券', 50000, '', 1, '2022-09-23 06:31:33', '2022-09-23 06:31:33', NULL);
INSERT INTO `coupon_info` VALUES (6, '满2千减5百优惠券', 50000, '', 1, '2022-09-23 06:33:21', '2022-09-23 06:33:21', NULL);
INSERT INTO `coupon_info` VALUES (7, '满2千减5百优惠券', 50000, '', 1, '2022-09-23 06:34:56', '2022-09-23 06:34:56', NULL);
INSERT INTO `coupon_info` VALUES (8, '满2千减5百优惠券', 50000, '', 1, '2022-09-23 06:36:17', '2022-09-23 06:36:17', NULL);
INSERT INTO `coupon_info` VALUES (9, '满2千减5百优惠券', 50000, '', 1, '2022-09-23 06:38:41', '2022-09-23 06:38:41', NULL);
INSERT INTO `coupon_info` VALUES (10, '满2千减5百优惠券', 50000, '0', 1, '2022-09-25 15:32:34', '2022-09-25 15:32:34', NULL);
INSERT INTO `coupon_info` VALUES (11, '满2千减5百优惠券', 50000, '0', 1, '2022-09-25 15:32:40', '2022-09-25 15:32:40', NULL);
INSERT INTO `coupon_info` VALUES (12, '满2千减5百优惠券', 50000, '0', 1, '2022-09-25 15:33:23', '2022-09-25 15:33:23', NULL);
INSERT INTO `coupon_info` VALUES (13, '满2千减5百优惠券', 50000, '0', 1, '2022-09-25 15:33:54', '2022-09-25 15:33:54', NULL);

-- ----------------------------
-- Table structure for file_info
-- ----------------------------
DROP TABLE IF EXISTS `file_info`;
CREATE TABLE `file_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '图片名称',
  `src` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '本地文件存储路径',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'URL地址',
  `user_id` int(11) NOT NULL DEFAULT 0 COMMENT '用户id',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of file_info
-- ----------------------------

-- ----------------------------
-- Table structure for goods_info
-- ----------------------------
DROP TABLE IF EXISTS `goods_info`;
CREATE TABLE `goods_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pic_url` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '图片',
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '商品名称',
  `price` int(11) NOT NULL DEFAULT 1 COMMENT '价格 单位分',
  `level1_category_id` int(11) NOT NULL COMMENT '1级分类id',
  `level2_category_id` int(11) NOT NULL DEFAULT 0 COMMENT '2级分类id',
  `level3_category_id` int(11) NOT NULL DEFAULT 0 COMMENT '3级分类id',
  `brand` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '品牌',
  `stock` int(11) NOT NULL DEFAULT 0 COMMENT '库存',
  `sale` int(11) NOT NULL DEFAULT 0 COMMENT '销量',
  `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标签',
  `detail_info` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '商品详情',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '商品表' ROW_FORMAT = Dynamic;


-- ----------------------------
-- Records of goods_info
-- ----------------------------
-- 家用电器 - 电视
INSERT INTO `goods_info` VALUES (1, 'https://img12.360buyimg.com/n7/jfs/t1/202536/20/30083/123456/632c9a29E1a7a5a5a/1a1a1a1a1a1a1a1a.jpg', '索尼XR-65X90J 65英寸 4K超高清 HDR 安卓智能电视', 699900, 1, 6, 20, '索尼', 100, 500, '4K超高清, HDR, 安卓智能, 全面屏', '索尼XR-65X90J采用XR认知芯片，4K HDR图像处理，全面屏设计，安卓智能系统，支持多种流媒体平台，音效震撼，是家庭影音娱乐的理想选择。', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- 家用电器 - 空调
INSERT INTO `goods_info` VALUES (2, 'https://img13.360buyimg.com/n7/jfs/t1/198818/30/30076/123456/632c9a29E1a7a5a5a/2b2b2b2b2b2b2b2b.jpg', '格力KFR-35GW/NhGc1B 1.5匹 变频冷暖 一级能效 壁挂式空调', 329900, 1, 7, 24, '格力', 200, 800, '变频冷暖, 一级能效, 壁挂式, 静音', '格力1.5匹变频空调，一级能效，冷暖两用，静音设计，适合15-20平米房间使用，节能省电，制冷制热效果好。', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- 手机数码 - 手机
INSERT INTO `goods_info` VALUES (3, 'https://img14.360buyimg.com/n7/jfs/t1/187655/33/29876/123456/632c9a29E1a7a5a5a/3c3c3c3c3c3c3c3c.jpg', 'Apple iPhone 14 Pro 256GB 暗紫色 移动联通电信5G手机', 899900, 2, 11, 28, '苹果', 150, 1200, '5G, 256GB, 暗紫色, A16芯片', 'iPhone 14 Pro采用A16仿生芯片，灵动岛设计，4800万像素主摄，支持ProMotion自适应刷新率，全天候显示，是苹果旗舰手机的代表。', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- 手机数码 - 耳机
INSERT INTO `goods_info` VALUES (4, 'https://img15.360buyimg.com/n7/jfs/t1/198818/30/30076/123456/632c9a29E1a7a5a5a/4d4d4d4d4d4d4d4d.jpg', '华为FreeBuds Pro 2 主动降噪 无线蓝牙耳机', 129900, 2, 12, 0, '华为', 300, 2000, '主动降噪, 无线蓝牙, 高清音质', '华为FreeBuds Pro 2采用主动降噪技术，高清音质，无线蓝牙连接，续航持久，佩戴舒适，适合日常使用和通勤。', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- 电脑办公 - 笔记本电脑
INSERT INTO `goods_info` VALUES (5, 'https://img16.360buyimg.com/n7/jfs/t1/187655/33/29876/123456/632c9a29E1a7a5a5a/5e5e5e5e5e5e5e5e.jpg', '联想ThinkPad X1 Carbon 2023款 14英寸轻薄笔记本电脑', 1099900, 3, 15, 32, '联想', 80, 300, '轻薄本, 14英寸, 高性能, 商务办公', '联想ThinkPad X1 Carbon 2023款采用14英寸轻薄设计，高性能处理器，长续航，商务办公首选，品质可靠，便携性强。', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- 家居家装 - 沙发
INSERT INTO `goods_info` VALUES (6, 'https://img17.360buyimg.com/n7/jfs/t1/198818/30/30076/123456/632c9a29E1a7a5a5a/6f6f6f6f6f6f6f6f.jpg', '顾家家居 现代简约布艺沙发 三人位沙发', 499900, 4, 0, 0, '顾家家居', 50, 150, '现代简约, 布艺沙发, 三人位, 客厅家具', '顾家家居现代简约布艺沙发，三人位设计，适合客厅使用，面料舒适，坐感柔软，简约时尚，易于清洁。', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- 服饰鞋包 - 运动鞋
INSERT INTO `goods_info` VALUES (7, 'https://img18.360buyimg.com/n7/jfs/t1/187655/33/29876/123456/632c9a29E1a7a5a5a/7g7g7g7g7g7g7g7g.jpg', 'Nike Air Zoom Pegasus 39 男子跑步鞋', 89900, 5, 0, 0, 'Nike', 120, 600, '跑步鞋, 气垫缓震, 轻便透气', 'Nike Air Zoom Pegasus 39男子跑步鞋，采用Zoom Air气垫缓震，轻便透气，适合日常跑步和训练，舒适耐穿，时尚美观。', '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
-- ----------------------------
-- Table structure for goods_options_info
-- ----------------------------
DROP TABLE IF EXISTS `goods_options_info`;
CREATE TABLE `goods_options_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `goods_id` int(11) NOT NULL COMMENT '商品id',
  `pic_url` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '图片',
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '商品名称',
  `price` int(11) NOT NULL DEFAULT 1 COMMENT '价格 单位分',
  `stock` int(11) NOT NULL DEFAULT 0 COMMENT '库存',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 27 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '商品规格表\n' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of goods_options_info
-- ----------------------------
INSERT INTO `goods_options_info` VALUES (1, 1, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '索尼XR-55X90J 55英寸 4K超高清', 549900, 50, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (2, 1, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '索尼XR-65X90J 65英寸 4K超高清', 699900, 30, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (3, 1, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '索尼XR-75X90J 75英寸 4K超高清', 899900, 20, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (4, 2, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '格力KFR-35GW 1.5匹 壁挂式空调', 329900, 100, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (5, 2, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '格力KFR-50GW 2匹 壁挂式空调', 429900, 80, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (6, 2, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '格力KFR-72LW 3匹 立柜式空调', 599900, 60, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (7, 3, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', 'iPhone 14 Pro 128GB 深空黑色', 799900, 80, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (8, 3, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', 'iPhone 14 Pro 256GB 暗紫色', 899900, 70, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (9, 3, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', 'iPhone 14 Pro 512GB 银色', 1099900, 50, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (10, 3, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', 'iPhone 14 Pro 1TB 金色', 1299900, 30, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (11, 4, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '华为FreeBuds Pro 2 陶瓷白', 129900, 150, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (12, 4, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '华为FreeBuds Pro 2 星际黑', 129900, 150, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (13, 4, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '华为FreeBuds Pro 2 冰霜银', 129900, 120, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (14, 5, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', 'ThinkPad X1 Carbon i5-1340P 16GB 512GB', 999900, 40, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (15, 5, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', 'ThinkPad X1 Carbon i7-1360P 16GB 1TB', 1299900, 35, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (16, 5, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', 'ThinkPad X1 Carbon i7-1360P 32GB 2TB', 1599900, 25, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (17, 6, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '顾家布艺沙发 三人位 米白色', 499900, 30, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (18, 6, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '顾家布艺沙发 三人位 深灰色', 499900, 30, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (19, 6, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '顾家布艺沙发 四人位 米白色', 599900, 20, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (20, 6, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '顾家布艺沙发 L型转角 米白色', 799900, 15, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (21, 7, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', 'Nike Air Zoom Pegasus 39 41码 黑白色', 89900, 60, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (22, 7, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', 'Nike Air Zoom Pegasus 39 42码 黑白色', 89900, 60, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (23, 7, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', 'Nike Air Zoom Pegasus 39 43码 黑白色', 89900, 50, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (24, 7, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', 'Nike Air Zoom Pegasus 39 44码 黑白色', 89900, 40, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (25, 7, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', 'Nike Air Zoom Pegasus 39 42码 蓝橙色', 89900, 55, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);
INSERT INTO `goods_options_info` VALUES (26, 7, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', 'Nike Air Zoom Pegasus 39 43码 蓝橙色', 89900, 45, '2023-01-01 00:00:00', '2023-01-01 00:00:00', NULL);

-- ----------------------------
-- Table structure for order_goods_info
-- ----------------------------
DROP TABLE IF EXISTS `order_goods_info`;
CREATE TABLE `order_goods_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '商品维度的订单表',
  `order_id` int(11) NOT NULL DEFAULT 0 COMMENT '关联的主订单表',
  `goods_id` int(11) NOT NULL DEFAULT 0 COMMENT '商品id',
  `goods_options_id` int(11) NULL DEFAULT 0 COMMENT '商品规格id sku id',
  `count` int(11) NOT NULL COMMENT '商品数量',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `price` int(11) NOT NULL DEFAULT 0 COMMENT '订单金额 单位分',
  `coupon_price` int(11) NOT NULL DEFAULT 0 COMMENT '优惠券金额 单位分',
  `actual_price` int(11) NOT NULL DEFAULT 0 COMMENT '实际支付金额 单位分',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 25 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '文章（种草）表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of order_goods_info
-- ----------------------------

-- ----------------------------
-- Table structure for order_info
-- ----------------------------
DROP TABLE IF EXISTS `order_info`;
CREATE TABLE `order_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `number` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '订单编号',
  `user_id` int(11) NOT NULL DEFAULT 0 COMMENT '用户id',
  `pay_type` tinyint(1) NOT NULL DEFAULT 0 COMMENT '支付方式 1微信 2支付宝 3云闪付',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `pay_at` datetime NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '支付时间',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '订单状态： 1待支付 2已支付待发货 3已发货 4已收货待评价 5已评价',
  `consignee_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收货人姓名',
  `consignee_phone` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收货人手机号',
  `consignee_address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收货人详细地址',
  `price` int(11) NOT NULL DEFAULT 0 COMMENT '订单金额 单位分',
  `coupon_price` int(11) NOT NULL DEFAULT 0 COMMENT '优惠券金额 单位分',
  `actual_price` int(11) NOT NULL DEFAULT 0 COMMENT '实际支付金额 单位分',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '文章（种草）表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of order_info
-- ----------------------------

-- ----------------------------
-- Table structure for permission_info
-- ----------------------------
DROP TABLE IF EXISTS `permission_info`;
CREATE TABLE `permission_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '权限名称',
  `path` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '路径',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_name`(`name` ASC) USING BTREE COMMENT '名称唯一索引'
) ENGINE = InnoDB AUTO_INCREMENT = 12 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of permission_info
-- ----------------------------
INSERT INTO `permission_info` VALUES (1, '轮播图管理', '/backend/rotation', '2022-09-25 15:03:01', '2025-12-04 08:10:34', NULL);
INSERT INTO `permission_info` VALUES (2, '文章管理', '/backend/article', '2022-12-26 19:51:44', '2022-12-26 19:52:29', NULL);
INSERT INTO `permission_info` VALUES (3, '评论管理', '/backend/comment', '2022-12-26 19:52:01', '2022-12-26 19:52:01', NULL);
INSERT INTO `permission_info` VALUES (4, '商品管理', '/backend/goods', '2025-12-04 15:23:03', '2025-12-04 15:23:04', NULL);
INSERT INTO `permission_info` VALUES (5, '分类管理', '/backend/category', '2025-12-04 15:27:24', '2025-12-04 15:38:31', NULL);
INSERT INTO `permission_info` VALUES (6, '订单管理', '/backend/order', '2025-12-04 15:27:26', '2025-12-04 15:38:34', NULL);
INSERT INTO `permission_info` VALUES (7, '收货地址管理', '/backend/consignee', '2025-12-04 15:27:29', '2025-12-04 15:38:36', NULL);

-- ----------------------------
-- Table structure for position_info
-- ----------------------------
DROP TABLE IF EXISTS `position_info`;
CREATE TABLE `position_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pic_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '图片链接',
  `goods_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '商品名称',
  `link` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '跳转链接',
  `sort` tinyint(4) NOT NULL DEFAULT 0 COMMENT '排序',
  `goods_id` int(11) NOT NULL DEFAULT 0 COMMENT '商品id',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of position_info
-- ----------------------------

-- ----------------------------
-- Table structure for praise_info
-- ----------------------------
DROP TABLE IF EXISTS `praise_info`;
CREATE TABLE `praise_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '点赞表',
  `user_id` int(11) NOT NULL,
  `type` tinyint(1) NOT NULL COMMENT '点赞类型 1商品 2文章',
  `object_id` int(11) NOT NULL DEFAULT 0 COMMENT '点赞对象id 方便后期扩展',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_index`(`user_id` ASC, `type` ASC, `object_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of praise_info
-- ----------------------------

-- ----------------------------
-- Table structure for refund_info
-- ----------------------------
DROP TABLE IF EXISTS `refund_info`;
CREATE TABLE `refund_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '售后退款表',
  `number` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '售后订单号',
  `order_id` int(11) NOT NULL COMMENT '订单id',
  `goods_id` int(11) NOT NULL DEFAULT 0 COMMENT '要售后的商品id\n',
  `reason` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '退款原因',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态 1待处理 2同意退款 3拒绝退款\n',
  `user_id` int(11) NOT NULL COMMENT '用户id',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of refund_info
-- ----------------------------

-- ----------------------------
-- Table structure for role_info
-- ----------------------------
DROP TABLE IF EXISTS `role_info`;
CREATE TABLE `role_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '角色名称',
  `desc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '描述',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_index`(`name` ASC) USING BTREE COMMENT '角色昵称唯一索引'
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role_info
-- ----------------------------
INSERT INTO `role_info` VALUES (1, '比超级低一级管理员', '拥有系统所有权限', '2025-12-02 08:21:58', '2025-12-04 08:08:58', NULL);
INSERT INTO `role_info` VALUES (2, '普通管理员', '拥有基础管理权限', '2022-12-21 10:43:33', '2022-12-21 10:43:33', NULL);

-- ----------------------------
-- Table structure for role_permission_info
-- ----------------------------
DROP TABLE IF EXISTS `role_permission_info`;
CREATE TABLE `role_permission_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL DEFAULT 0 COMMENT '角色id',
  `permission_id` int(11) NOT NULL COMMENT '权限id',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_index`(`role_id` ASC, `permission_id` ASC) USING BTREE COMMENT '唯一索引'
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role_permission_info
-- ----------------------------
INSERT INTO `role_permission_info` VALUES (4, 2, 2, '2025-12-03 15:54:08', '2025-12-04 16:55:24');
INSERT INTO `role_permission_info` VALUES (5, 2, 3, '2025-12-03 15:54:10', '2025-12-04 16:55:27');
INSERT INTO `role_permission_info` VALUES (6, 1, 1, '2025-12-04 08:10:16', '2025-12-04 08:10:16');
INSERT INTO `role_permission_info` VALUES (7, 1, 2, '2025-12-04 08:10:16', '2025-12-04 08:10:16');
INSERT INTO `role_permission_info` VALUES (8, 1, 3, '2025-12-04 08:10:16', '2025-12-04 08:10:16');
INSERT INTO `role_permission_info` VALUES (9, 1, 4, '2025-12-04 08:10:16', '2025-12-04 08:10:16');
INSERT INTO `role_permission_info` VALUES (10, 1, 5, '2025-12-04 08:10:16', '2025-12-04 08:10:16');
INSERT INTO `role_permission_info` VALUES (11, 1, 6, '2025-12-04 08:10:16', '2025-12-04 08:10:16');
INSERT INTO `role_permission_info` VALUES (12, 1, 7, '2025-12-04 08:10:16', '2025-12-04 08:10:16');

-- ----------------------------
-- Table structure for rotation_info
-- ----------------------------
DROP TABLE IF EXISTS `rotation_info`;
CREATE TABLE `rotation_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pic_url` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '轮播图片',
  `link` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '跳转链接',
  `sort` tinyint(1) NOT NULL DEFAULT 0 COMMENT '排序字段',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '轮播图表\n' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of rotation_info
-- ----------------------------
INSERT INTO `rotation_info` VALUES (9, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '1', 1, '2025-12-02 08:34:08', '2025-12-03 07:09:15', NULL);
INSERT INTO `rotation_info` VALUES (14, 'https://daxinggonghui.oss-cn-beijing.aliyuncs.com/images/zhanwei.jpg', '2', 0, '2025-12-04 05:59:04', '2025-12-04 05:59:23', NULL);

-- ----------------------------
-- Table structure for seckill_goods
-- ----------------------------
DROP TABLE IF EXISTS `seckill_goods`;
CREATE TABLE `seckill_goods`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `goods_id` bigint(20) NOT NULL COMMENT '商品ID',
  `goods_options_id` bigint(20) NOT NULL COMMENT '商品规格ID',
  `original_price` int(11) NOT NULL COMMENT '原始价格 单位分',
  `seckill_price` int(11) NOT NULL COMMENT '秒杀价格 单位分',
  `seckill_stock` int(11) NOT NULL COMMENT '秒杀库存',
  `start_time` datetime NOT NULL COMMENT '秒杀开始时间',
  `end_time` datetime NOT NULL COMMENT '秒杀结束时间',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态：0-未开始 1-进行中 2-已结束',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_goods_id_options_id`(`goods_id` ASC, `goods_options_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_start_end_time`(`start_time` ASC, `end_time` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '秒杀商品表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of seckill_goods
-- ----------------------------
INSERT INTO `seckill_goods` VALUES (1, 1, 2, 699900, 599900, 10, '2025-12-01 00:00:00', '2025-12-31 23:59:59', 1, '2025-11-30 10:00:00', '2025-11-30 10:00:00');
INSERT INTO `seckill_goods` VALUES (2, 3, 8, 899900, 799900, 5, '2025-12-01 00:00:00', '2025-12-31 23:59:59', 1, '2025-11-30 10:00:00', '2025-11-30 10:00:00');
INSERT INTO `seckill_goods` VALUES (3, 4, 11, 129900, 99900, 20, '2025-12-01 00:00:00', '2025-12-31 23:59:59', 1, '2025-11-30 10:00:00', '2025-11-30 10:00:00');

-- ----------------------------
-- Table structure for seckill_order
-- ----------------------------
DROP TABLE IF EXISTS `seckill_order`;
CREATE TABLE `seckill_order`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `order_id` bigint(20) NOT NULL COMMENT '订单ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `goods_id` bigint(20) NOT NULL COMMENT '商品ID',
  `goods_options_id` bigint(20) NOT NULL COMMENT '商品规格ID',
  `seckill_price` int(11) NOT NULL COMMENT '秒杀价格 单位分',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态：0-待支付 1-已支付 2-已取消',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`, `user_id`) USING BTREE,
  UNIQUE INDEX `uk_order_user`(`order_id` ASC, `user_id` ASC) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_goods_id_options_id`(`goods_id` ASC, `goods_options_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '秒杀订单表' ROW_FORMAT = DYNAMIC PARTITION BY KEY (`user_id`)
PARTITIONS 16
(PARTITION `p0` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p1` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p10` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p11` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p12` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p13` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p14` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p15` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p2` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p3` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p4` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p5` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p6` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p7` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p8` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p9` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 )
;

-- ----------------------------
-- Records of seckill_order
-- ----------------------------

-- ----------------------------
-- Table structure for user_coupon_info
-- ----------------------------
DROP TABLE IF EXISTS `user_coupon_info`;
CREATE TABLE `user_coupon_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户优惠券表',
  `user_id` int(11) NOT NULL DEFAULT 0,
  `coupon_id` int(11) NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1可用 2已用 3过期',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_coupon_info
-- ----------------------------
INSERT INTO `user_coupon_info` VALUES (1, 1, 1, 1, '2022-07-29 16:01:13', '2022-07-29 16:01:13', NULL);
INSERT INTO `user_coupon_info` VALUES (2, 1, 2, 1, '2022-07-29 16:01:13', '2022-07-29 16:01:13', NULL);

-- ----------------------------
-- Table structure for user_info
-- ----------------------------
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `avatar` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '头像',
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `user_salt` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '加密盐 生成密码用',
  `sex` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1男 2女',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1正常 2拉黑冻结',
  `sign` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '个性签名',
  `secret_answer` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密保问题的答案',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '商品表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_info
-- ----------------------------
INSERT INTO `user_info` VALUES (1, 'wangzhongyang', 'http://dummyimage.com/100x100', '82131d93ab13a1a4f9ec840a9ddbabf7', 'T0iKtv31BU', 1, 1, '和我一起学编程吧', '六个1', '2024-12-26 11:25:43', '2025-12-04 07:56:57', NULL);

SET FOREIGN_KEY_CHECKS = 1;
