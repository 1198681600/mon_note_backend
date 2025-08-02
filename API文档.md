# MoodTrace 后端 API 文档

## 基本信息

- **服务地址**: http://localhost:8080
- **API 前缀**: `/api`
- **认证方式**: Bearer Token (请求头: `Authorization: Bearer <token>`)

## 响应格式

所有 API 响应都采用统一格式：

```json
{
    "code": 200,
    "message": "操作成功",
    "data": {}
}
```

## 认证相关 API

### 1. 发送邮箱验证码

**接口地址**: `POST /api/auth/send-code`

**请求参数**:
```json
{
    "email": "user@example.com"
}
```

**响应示例**:
```json
{
    "code": 200,
    "message": "验证码已发送"
}
```

**说明**: 验证码固定为 `111111`，有效期15分钟

---

### 2. 用户注册

**接口地址**: `POST /api/auth/register`

**请求参数**:
```json
{
    "email": "user@example.com",
    "code": "111111"
}
```

**参数说明**:
- `email`: 邮箱地址，必须是有效邮箱格式
- `code`: 邮箱验证码

**响应示例**:
```json
{
    "code": 200,
    "message": "注册成功"
}
```

---

### 3. 验证邮箱

**接口地址**: `POST /api/auth/verify-email`

**请求参数**:
```json
{
    "email": "user@example.com",
    "code": "111111"
}
```

**响应示例**:
```json
{
    "code": 200,
    "message": "邮箱验证成功"
}
```

---

### 4. 用户登录

**接口地址**: `POST /api/auth/login`

**请求参数**:
```json
{
    "email": "user@example.com",
    "code": "111111"
}
```

**参数说明**:
- `email`: 邮箱地址
- `code`: 邮箱验证码

**响应示例**:
```json
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "token": "a1b2c3d4e5f6...",
        "user": {
            "id": 1,
            "email": "user@example.com",
            "is_email_verified": true,
            "avatar": "",
            "nickname": "",
            "gender": "",
            "age": 0,
            "profession": "",
            "created_at": "2025-01-01T12:00:00Z",
            "updated_at": "2025-01-01T12:00:00Z"
        }
    }
}
```

**说明**: 
- 使用邮箱验证码登录，无需密码
- Token 有效期24小时
- 请将 token 保存并在后续请求中使用

---

## 情绪分析 API

### 1. 分析日记内容情绪

**接口地址**: `POST /api/emotion/analyze-diary`

**请求头**:
```
Authorization: Bearer <your_token>
```

**请求参数**:
```json
{
    "diary_content": "今天上午心情很好，工作很顺利。中午吃饭的时候和同事聊天很开心。下午开会的时候有点紧张，但是最后顺利完成了。晚上回家看到家人很温暖。",
    "diary_date": "2025-01-15",
    "user_context": {
        "age": 25,
        "gender": "女",
        "profession": "软件工程师"
    }
}
```

**参数说明**:
- `diary_content`: 日记内容 (必填)
- `diary_date`: 日记日期 YYYY-MM-DD格式 (必填)
- `user_context`: 用户背景信息 (可选，有助于更精准的分析)
    - `age`: 年龄
    - `gender`: 性别  
    - `profession`: 职业

**响应示例**:
```json
{
    "code": 200,
    "message": "情绪分析完成",
    "data": {
        "emotions": [
            {
                "emotion": "开心",
                "intensity": 0.8,
                "time_period": "上午",
                "color": "#FFD700",
                "description": "工作顺利带来的满足感"
            },
            {
                "emotion": "紧张", 
                "intensity": 0.6,
                "time_period": "下午",
                "color": "#E17055",
                "description": "开会前的焦虑情绪"
            },
            {
                "emotion": "温暖",
                "intensity": 0.9,
                "time_period": "晚上", 
                "color": "#FFB8B8",
                "description": "家庭温馨带来的感动"
            }
        ],
        "gradient_suggestion": {
            "type": "sweep",
            "reasoning": "一天中有多种情绪转换，适合使用扫描渐变展现情绪轮转"
        },
        "summary": {
            "dominant_emotion": "积极",
            "emotional_stability": 7.5,
            "mood_trend": "整体向好",
            "energy_level": "中等偏高"
        },
        "insights": [
            "今天整体情绪状态良好，工作和家庭都给你带来了正面情绪",
            "下午的紧张情绪是正常的工作压力反应，但最终得到了很好的缓解",
            "家庭关系是你重要的情绪支撑点"
        ],
        "recommendations": [
            "继续保持工作和生活的平衡",
            "可以在开会前做一些深呼吸来缓解紧张",
            "多和家人分享工作中的成就"
        ]
    }
}
```

---

### 2. 分析一周情绪趋势

**接口地址**: `POST /api/emotion/analyze-weekly`

**请求头**:
```
Authorization: Bearer <your_token>
```

**请求参数**:
```json
{
    "week_start": "2025-01-13",
    "diary_data": [
        {
            "date": "2025-01-13", 
            "emotions": [
                {
                    "emotion": "开心",
                    "intensity": 0.8,
                    "time_period": "上午"
                }
            ]
        },
        {
            "date": "2025-01-14",
            "emotions": [
                {
                    "emotion": "焦虑", 
                    "intensity": 0.7,
                    "time_period": "下午"
                }
            ]
        }
    ]
}
```

**参数说明**:
- `week_start`: 一周开始日期 YYYY-MM-DD格式 (必填)
- `diary_data`: 这一周的情绪数据数组 (必填)

**响应示例**:
```json
{
    "code": 200,
    "message": "一周情绪分析完成", 
    "data": {
        "weekly_pattern": {
            "monday": {"dominant": "平静", "intensity": 0.6},
            "tuesday": {"dominant": "开心", "intensity": 0.8},
            "wednesday": {"dominant": "焦虑", "intensity": 0.7},
            "thursday": {"dominant": "满足", "intensity": 0.7},
            "friday": {"dominant": "兴奋", "intensity": 0.9},
            "saturday": {"dominant": "放松", "intensity": 0.8},
            "sunday": {"dominant": "温暖", "intensity": 0.8}
        },
        "insights": [
            "本周整体情绪波动属于正常范围",
            "周三出现焦虑情绪，可能与工作压力相关",
            "周末情绪明显改善，休息对你很重要"
        ],
        "recommendations": [
            "周三是你的情绪低谷期，建议安排轻松的活动",
            "保持现有的周末放松习惯",
            "可以在周中增加一些减压活动"
        ],
        "emotion_score": 7.4,
        "stability_score": 6.8
    }
}
```

**说明**:
- 需要用户登录认证
- 基于Claude AI进行智能分析
- 返回详细的一周情绪模式和建议

---

### 3. Claude AI 提示词模板

**服务端实现参考** - 调用Claude API时使用以下提示词：

```
你是一个专业的情绪分析专家，专门分析用户的日记内容并提取情绪信息。

## 任务要求
分析以下日记内容，提取其中的情绪信息，并按指定JSON格式返回结果。

## 日记内容
{diary_content}

## 日记日期  
{diary_date}

## 用户背景（可选）
年龄: {age}
性别: {gender}  
职业: {profession}

## 情绪色彩映射表
请从以下预定义情绪中选择最匹配的：
- 开心: #FFD700 (金黄色)
- 快乐: #FF6B6B (珊瑚红)
- 兴奋: #FF8E53 (橙色)
- 平静: #4ECDC4 (青绿色)
- 放松: #45B7D1 (天蓝色)
- 满足: #96CEB4 (薄荷绿)
- 难过: #6C5CE7 (紫色)
- 沮丧: #74B9FF (蓝色)
- 焦虑: #FDCB6E (淡黄色)
- 紧张: #E17055 (橙红色)
- 愤怒: #D63031 (红色)
- 烦躁: #E84393 (粉红色)
- 感动: #A29BFE (淡紫色)
- 温暖: #FFB8B8 (粉色)
- 孤独: #636E72 (灰色)
- 迷茫: #B2BEC3 (浅灰色)

## 输出格式
请严格按照以下JSON格式返回，不要包含任何其他文字：

{
  "emotions": [
    {
      "emotion": "情绪名称",
      "intensity": 0.0-1.0之间的数值,
      "time_period": "上午/中午/下午/晚上",
      "color": "对应的颜色代码",
      "description": "该情绪的简短描述(20字以内)"
    }
  ],
  "gradient_suggestion": {
    "type": "radial/sweep/multiPoint/wave/diagonal",
    "reasoning": "选择此渐变类型的原因"
  },
  "summary": {
    "dominant_emotion": "主导情绪",
    "emotional_stability": 1-10的评分,
    "mood_trend": "整体趋势描述",
    "energy_level": "低/中等/高"
  },
  "insights": [
    "分析洞察1",
    "分析洞察2",
    "分析洞察3"
  ],
  "recommendations": [
    "建议1", 
    "建议2",
    "建议3"
  ]
}

## 注意事项
1. intensity值必须在0.0-1.0之间
2. emotion必须从预定义列表中选择
3. time_period只能是：上午、中午、下午、晚上
4. 每个数组最多包含5个元素
5. 描述要简洁准确，贴合中文表达习惯
6. 分析要专业但易懂，避免过于学术化的表述
```

**说明**:
- 需要用户登录认证
- 基于Claude AI进行智能分析
- 返回详细的情绪数据用于客户端渐变背景生成

---

## 需要认证的 API

以下接口需要在请求头中携带 `Authorization: Bearer <token>`

### 2. 获取用户信息

**接口地址**: `GET /api/profile`

**请求头**:
```
Authorization: Bearer <your_token>
```

**响应示例**:
```json
{
    "code": 200,
    "message": "获取用户信息成功",
    "data": {
        "id": 1,
        "email": "user@example.com",
        "is_email_verified": true,
        "avatar": "",
        "nickname": "",
        "gender": "",
        "age": 0,
        "profession": "",
        "created_at": "2025-01-01T12:00:00Z",
        "updated_at": "2025-01-01T12:00:00Z"
    }
}
```

---

### 3. 更新用户信息

**接口地址**: `POST /api/profile`

**请求头**:
```
Authorization: Bearer <your_token>
```

**请求参数**:
```json
{
    "nickname": "我的昵称",
    "gender": "男",
    "age": 25,
    "profession": "软件工程师",
    "avatar": "https://example.com/avatar.jpg"
}
```

**参数说明**:
- `nickname`: 昵称 (可选)
- `gender`: 性别 (可选)
- `age`: 年龄 (可选)
- `profession`: 职业 (可选)
- `avatar`: 头像URL (可选)

**响应示例**:
```json
{
    "code": 200,
    "message": "更新用户信息成功",
    "data": {
        "id": 1,
        "email": "user@example.com",
        "is_email_verified": true,
        "avatar": "https://example.com/avatar.jpg",
        "nickname": "我的昵称",
        "gender": "男",
        "age": 25,
        "profession": "软件工程师",
        "created_at": "2025-01-01T12:00:00Z",
        "updated_at": "2025-01-01T14:30:00Z"
    }
}
```

**说明**: 
- 所有参数都是可选的，只更新提供的字段
- 需要用户登录认证

---

### 4. 生成文件上传链接

**接口地址**: `POST /api/upload/generate-url`

**请求头**:
```
Authorization: Bearer <your_token>
```

**请求参数**:
```json
{
    "file_type": "image/jpeg"
}
```

**参数说明**:
- `file_type`: 文件类型，支持的类型：
  - `image/jpeg` 或 `image/jpg`: JPEG图片
  - `image/png`: PNG图片
  - `image/webp`: WebP图片

**响应示例**:
```json
{
    "code": 200,
    "message": "生成上传链接成功",
    "data": {
        "upload_url": "https://account_id.r2.cloudflarestorage.com/bucket/avatars/uuid.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&...",
        "file_url": "https://account_id.r2.cloudflarestorage.com/bucket/avatars/uuid.jpg",
        "file_name": "avatars/uuid.jpg",
        "expires_in": 900
    }
}
```

**使用流程**:
1. 调用此接口获取预签名上传URL
2. 使用 `upload_url` 通过 HTTP PUT 请求直接上传文件到R2
3. 上传成功后，使用 `file_url` 作为头像URL调用更新用户信息接口

**上传文件示例**:
```bash
# 1. 获取上传链接
curl -X POST http://localhost:8080/api/upload/generate-url \
  -H "Authorization: Bearer your_token_here" \
  -H "Content-Type: application/json" \
  -d '{"file_type":"image/jpeg"}'

# 2. 使用返回的upload_url上传文件
curl -X PUT "返回的upload_url" \
  -H "Content-Type: image/jpeg" \
  --data-binary @avatar.jpg

# 3. 更新用户头像
curl -X POST http://localhost:8080/api/profile \
  -H "Authorization: Bearer your_token_here" \
  -H "Content-Type: application/json" \
  -d '{"avatar":"返回的file_url"}'
```

**说明**: 
- 上传链接有效期15分钟
- 文件会保存到 `avatars/` 目录下
- 需要用户登录认证
- 需要配置Cloudflare R2环境变量

---

### 5. 退出登录

**接口地址**: `POST /api/logout`

**请求头**:
```
Authorization: Bearer <your_token>
```

**响应示例**:
```json
{
    "code": 200,
    "message": "退出登录成功"
}
```

**说明**: 退出登录后 token 将失效

---

## 日记 API

### 1. 创建日记

**接口地址**: `POST /api/diary/create`

**说明**: 创建一篇新的日记。每个用户每天只能创建一篇日记。

**请求头**:
`Authorization: Bearer <your_token>`

**查询参数**:
- `date`: 日期，格式为 `YYYY-MM-DD` (可选, 默认是当天的 UTC+8 日期)。

**请求参数**:
```json
{
    "content": "这是我的第一篇日记。"
}
```

**参数说明**:
- `content`: 日记内容 (必填)

**响应示例**:
```json
{
    "ID": 1,
    "CreatedAt": "2025-08-03T10:00:00Z",
    "UpdatedAt": "2025-08-03T10:00:00Z",
    "DeletedAt": null,
    "UserID": 1,
    "Date": "2025-08-03T00:00:00Z",
    "Content": "这是我的第一篇日记。"
}
```

**错误响应**:
```json
{
    "error": "diary already exists for this date"
}
```

---

### 2. 获取用户的所有日记

**接口地址**: `POST /api/diary/list`

**说明**: 获取当前用户的所有日记，按创建时间降序排列。

**请求头**:
`Authorization: Bearer <your_token>`

**请求参数**: 无需请求体

**响应示例**:
```json
[
    {
        "ID": 2,
        "CreatedAt": "2025-08-04T10:00:00Z",
        "UpdatedAt": "2025-08-04T10:00:00Z",
        "DeletedAt": null,
        "UserID": 1,
        "Date": "2025-08-04T00:00:00Z",
        "Content": "这是我的第二篇日记。"
    },
    {
        "ID": 1,
        "CreatedAt": "2025-08-03T10:00:00Z",
        "UpdatedAt": "2025-08-03T10:00:00Z",
        "DeletedAt": null,
        "UserID": 1,
        "Date": "2025-08-03T00:00:00Z",
        "Content": "这是我的第一篇日记。"
    }
]
```

---

### 3. 获取单篇日记

**接口地址**: `POST /api/diary/get`

**说明**: 获取指定 ID 的日记。

**请求头**:
`Authorization: Bearer <your_token>`

**请求参数**:
```json
{
    "id": 1
}
```

**参数说明**:
- `id`: 日记的唯一 ID (必填)

**响应示例**:
```json
{
    "ID": 1,
    "CreatedAt": "2025-08-03T10:00:00Z",
    "UpdatedAt": "2025-08-03T10:00:00Z",
    "DeletedAt": null,
    "UserID": 1,
    "Date": "2025-08-03T00:00:00Z",
    "Content": "这是我的第一篇日记。"
}
```

**错误响应**:
```json
{
    "error": "unauthorized"
}
```
或
```json
{
    "error": "Diary not found"
}
```

---

### 4. 更新日记

**接口地址**: `POST /api/diary/update`

**说明**: 更新指定 ID 的日记。

**请求头**:
`Authorization: Bearer <your_token>`

**请求参数**:
```json
{
    "id": 1,
    "content": "更新后的日记内容。"
}
```

**参数说明**:
- `id`: 日记的唯一 ID (必填)
- `content`: 更新后的日记内容 (必填)

**响应示例**:
```json
{
    "ID": 1,
    "CreatedAt": "2025-08-03T10:00:00Z",
    "UpdatedAt": "2025-08-03T11:00:00Z",
    "DeletedAt": null,
    "UserID": 1,
    "Date": "2025-08-03T00:00:00Z",
    "Content": "更新后的日记内容。"
}
```

**错误响应**:
```json
{
    "error": "unauthorized"
}
```

---

### 5. 删除日记

**接口地址**: `POST /api/diary/delete`

**说明**: 删除指定 ID 的日记。

**请求头**:
`Authorization: Bearer <your_token>`

**请求参数**:
```json
{
    "id": 1
}
```

**参数说明**:
- `id`: 日记的唯一 ID (必填)

**响应示例**:
```json
{
    "message": "Diary deleted successfully"
}
```

**错误响应**:
```json
{
    "error": "unauthorized"
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 操作成功 |
| 400 | 请求参数错误 |
| 401 | 未授权/token无效 |
| 403 | 无权限访问 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 500 | 服务器内部错误 |

## 常见错误响应

### 参数错误
```json
{
    "code": 400,
    "message": "请求参数错误"
}
```

### 用户已存在
```json
{
    "code": 400,
    "message": "用户已存在"
}
```

### 验证码错误
```json
{
    "code": 400,
    "message": "验证码无效或已过期"
}
```

### 用户不存在
```json
{
    "code": 400,
    "message": "用户不存在"
}
```

### 邮箱未验证
```json
{
    "code": 400,
    "message": "请先验证邮箱"
}
```


### Token无效
```json
{
    "code": 401,
    "message": "token无效或已过期"
}
```

## 使用流程示例

### 完整注册登录流程

#### 注册流程
1. **发送验证码**
   ```bash
   curl -X POST http://localhost:8080/api/auth/send-code \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com"}'
   ```

2. **用户注册**（使用验证码）
   ```bash
   curl -X POST http://localhost:8080/api/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","code":"111111"}'
   ```

#### 登录流程
1. **发送验证码**
   ```bash
   curl -X POST http://localhost:8080/api/auth/send-code \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com"}'
   ```

2. **用户登录**（使用验证码）
   ```bash
   curl -X POST http://localhost:8080/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","code":"111111"}'
   ```

3. **获取用户信息**
   ```bash
   curl -X GET http://localhost:8080/api/profile \
     -H "Authorization: Bearer your_token_here"
   ```

4. **更新用户信息**
   ```bash
   curl -X POST http://localhost:8080/api/profile \
     -H "Authorization: Bearer your_token_here" \
     -H "Content-Type: application/json" \
     -d '{
       "nickname": "我的昵称",
       "gender": "男",
       "age": 25,
       "profession": "软件工程师",
       "avatar": "https://example.com/avatar.jpg"
     }'
   ```

5. **退出登录**
   ```bash
   curl -X POST http://localhost:8080/api/logout \
     -H "Authorization: Bearer your_token_here"
   ```

6. **分析日记情绪**
   ```bash
   curl -X POST http://localhost:8080/api/emotion/analyze-diary \
     -H "Authorization: Bearer your_token_here" \
     -H "Content-Type: application/json" \
     -d '{
       "diary_content": "今天心情很好，工作很顺利。",
       "diary_date": "2025-01-15",
       "user_context": {
         "age": 25,
         "gender": "女",
         "profession": "软件工程师"
       }
     }'
   ```

## 数据库表结构

### Users 表
- `id`: 用户ID (主键)
- `email`: 邮箱地址 (唯一索引)
- `is_email_verified`: 邮箱是否已验证
- `avatar`: 头像URL
- `nickname`: 昵称
- `gender`: 性别
- `age`: 年龄
- `profession`: 职业
- `created_at`: 创建时间
- `updated_at`: 更新时间

### Email Verifications 表
- `id`: 记录ID (主键)
- `email`: 邮箱地址
- `verification_code`: 验证码
- `is_used`: 是否已使用
- `expires_at`: 过期时间
- `created_at`: 创建时间

### User Sessions 表
- `id`: 会话ID (主键)
- `user_id`: 用户ID (外键)
- `token`: 会话token
- `expires_at`: 过期时间
- `created_at`: 创建时间
