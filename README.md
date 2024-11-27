# Ghost API

The Ghost API is designed to generate structured test responses dynamically. It supports various data types such as strings, numbers, dates, arrays, and objects. This documentation provides an overview of endpoints, supported data types, and response formats.

# Ghost API Documentation

create a `.ghostrc` file in the root of your project.
alternatively, you can pass the config file path as a flag `-c <path>`

**Config Example:**

```yaml
port: 8080 # Port to listen on
latency: 150 # Default latency in milliseconds(optional)
jitter: 50 # Default jitter in milliseconds(optional) (the range is [-jitter, +jitter])
endpoints: # Array of endpoints
  - name: "GetUser" # Name of the endpoint
    url: "/api/events" # URL of the endpoint
    response:
      status_code: 200
      data_type: "application/json"
      data:
        type: array
        range: [6, 20]
        items:
          type: object
          fields:
            id:
              type: number
              range: [1, 100]
            name:
              type: string
              metadata: name
              range: [6, 100]
            imgUrl:
              type: string
              value: "https://picsum.photos/200/300"
            date:
              type: date
              range: ["3d", "60d"]
```

## Types of Data

```yaml
string:
  type: string
  value: ~ # fixed value (optional)
  metadata: ~ # Specifies type (e.g., name, first_name, email); defaults to lorem ipsum
  range: [5, 20] # Length range [min, max](optional)

number:
  type: number
  value: ~ # fixed value (optional)
  metadata: ~ # Specifies type (e.g., age, latitude); defaults to random value
  range: [1, 100] # Value range [min, max]

date:
  type: date
  value: ~ # Default value (optional)
  range: [-30d, 30m] # Range in days from the present [min_offset, max_offset]
  format: "2006-01-02" # Date format (Go-style)

array:
  type: array
  range: [1, 5] # Number of items [min, max]
  items: # Schema of array elements
    type: string
    metadata: email

object:
  type: object
  value: ~ # Default object (optional)
  fields: # Key-value schemas
    title:
      type: string
      range: [5, 50]
    age:
      type: numeric
      range: [18, 65]
    created_at:
      type: date
      format: "2006-01-02"
```

**Response Example:**

```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "age": 29,
  "created_at": "2023-11-27",
  "preferences": "dark_mode"
}
```

---

## **Supported Data Types**

### **1. String**

- **Properties:**
  - `value`: Default value (optional).
  - `metadata`: Specifies type (`name`, `fullname`, `email`).
  - `range`: Length range `[min, max]`.

#### Example:

```yaml
type: string,
metadata: name,
range: [5, 10]
```

### **3. Number**

- **Properties:**
  - `value`: Default value (optional).
  - `metadata`: Specifies type (`age`, `latitude`, etc.).
  - `range`: Value range `[min, max]`.

#### Example:

```yaml
type: numeric,
metadata: age,
range: [18, 65]
```

### **4. Date**

- **Properties:**
  - `value`: Default value (optional).
  - `range`: Range from the current date `[min_offset, max_offset]` in days.
  - `format`: Date format (uses Go’s format).

#### Example:

```yaml
type: date,
format: "2006-01-02",
range: ["-30d", "30"]
```

### **5. Array**

- **Properties:**
  - `range`: Number of items `[min, max]`.
  - `items`: Schema of array elements.

#### Example:

```yaml
type: array,
range: [1, 5],
items:
  type: string,
  metadata: email
```

### **6. Object**

- **Properties:**
  - `fields`: Nested key-value schemas of other types.

# Acknowledgements

This project could not have been possible without the following open source projects:
[go-faker](github.com/go-faker/faker/)
