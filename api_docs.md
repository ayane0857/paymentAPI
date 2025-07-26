## /

content で Hello World!!が返ってきます

## GET /payment

これまでの支払い履歴を閲覧できます

#### クエリパラメータ

| パラメータ名 | 型      | 必須 | 説明                                        |
| :----------- | :------ | :--- | :------------------------------------------ |
| `limit`      | integer | 任意 | 取得する支払い履歴の一覧。デフォルトは 10。 |
| `offset`     | integer | 任意 | オフセット。                                |

#### レスポンス

```json
[
  {
    "ID": 1,
    "Money": 110,
    "Payment": "PayPay",
    "CreatedAt": "2025-07-22T02:59:28.53506Z",
    "UpdatedAt": "2025-07-22T02:59:28.53506Z"
  },
  {
    "ID": 2,
    "Money": 220,
    "Payment": "PayPay",
    "CreatedAt": "2025-07-22T05:41:13.160653Z",
    "UpdatedAt": "2025-07-22T05:41:13.160653Z"
  }
]
```

## GET /payment/:id

特定の支払い履歴を閲覧できます

#### レスポンス

```json
{
  "ID": 1,
  "Money": 110,
  "Payment": "PayPay",
  "CreatedAt": "2025-07-22T02:59:28.53506Z",
  "UpdatedAt": "2025-07-22T02:59:28.53506Z"
}
```

## POST /payment

支払い履歴を追加します

```json
{
  "ID": 2,
  "Money": 220,
  "Payment": "PayPay",
  "CreatedAt": "2025-07-22T05:41:13.160653Z",
  "UpdatedAt": "2025-07-22T05:41:13.160653Z"
}
```

## PUT /payment/:id

支払い履歴を編集します

```json
{
  "ID": 1,
  "Money": 220,
  "Payment": "PayPay",
  "CreatedAt": "2025-07-22T02:59:28.53506Z",
  "UpdatedAt": "2025-07-22T05:46:21.484355757Z"
}
```

## DELETE /payment/:id

支払い履歴を削除します

```json
{
  "ID": 1,
  "Money": 220,
  "Payment": "PayPay",
  "CreatedAt": "2025-07-22T02:59:28.53506Z",
  "UpdatedAt": "2025-07-22T05:46:21.484355Z"
}
```

## GET /balance

現在の残高状況を表示します

## PUT /balance

現在の残高状況を更新します

```

```
