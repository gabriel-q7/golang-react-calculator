# API Reference

Base path: `/api`

All responses are JSON.

## `GET /api/health`

Liveness check.

**Response `200`**

```json
{ "status": "ok" }
```

## `POST /api/calculate`

Performs one arithmetic operation.

**Request body**

```json
{
  "a": 10,
  "b": 2,
  "operator": "divide"
}
```

| Field    | Type   | Description                                              |
| -------- | ------ | ---------------------------------------------------------- |
| a        | number | First operand                                              |
| b        | number | Second operand                                              |
| operator | string | One of `add`, `subtract`, `multiply`, `divide`              |

**Response `200`**

```json
{ "result": 5 }
```

**Response `400` — bad input**

Malformed JSON body, or a `divide` with `b == 0`:

```json
{ "error": "division by zero" }
```

**Response `422` — unsupported operator**

```json
{ "error": "unsupported operator: \"mod\"" }
```

**Response `405` — wrong HTTP method**

```json
{ "error": "method not allowed" }
```
