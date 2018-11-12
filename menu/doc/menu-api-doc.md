# Menu API document

### 1. Create a new menu item

**Request**

```
POST /menu/item
Content-Type: application/json
```

**Parameters**

| Parameter       | Type   | Description       |
| --------------- | ------ | ----------------- |
| itemName        | String | Name of the item  |
| itemPrice       | Float  | Price of the item |
| itemDescription | String | Ingredients used  |
| itemCalories    | String | Calorie content   |

*Unique item id will be generated by MongoDB*

**Response**

Parameters for Success (Status code: 200)

| Parameter       | Type   | Description       |
| --------------- | ------ | ----------------- |
| itemId          | string | Unique item Id    |
| itemName        | String | Name of the item  |
| itemPrice       | Float  | Price of the item |
| itemDescription | String | Ingredients used  |
| itemCalories    | String | Calorie content   |

Parameters for Error (Status code: 401)

| Parameter | Type   | Description                   |
| --------- | ------ | ----------------------------- |
| message   | string | Not authorized to create item |



### 2. Edit a menu item

**Request**

```
PUT /menu/item/{id}
Content-Type: application/json
```

**Parameters**

| Parameter       | Type   | Description       |
| --------------- | ------ | ----------------- |
| itemName        | string | Name of the item  |
| itemPrice       | float  | Price of the item |
| itemDescription | string | Ingredients used  |
| itemCalories    | float  | Calorie content   |

**Response**

Parameters for Success (Status code: 200)

| Parameter       | Type   | Description       |
| --------------- | ------ | ----------------- |
| itemId          | string | Unique item Id    |
| itemName        | float  | Name of the item  |
| itemPrice       | string | Price of the item |
| itemDescription | string | Ingredients used  |
| itemCalories    | float  | Calorie content   |

Parameters for Error (Status code: 401)

| Parameter | Type   | Description                   |
| --------- | ------ | ----------------------------- |
| message   | string | Not authorized to update item |



### 3. Get a menu item based on item ID

**Request**

```
GET /menu/item/:id
Content-Type: application/json
```

Parameters

| Parameter | Type   | Description                                                  |
| --------- | ------ | ------------------------------------------------------------ |
| id        | string | Unique item id  generated by MongoDB, for example: ObjectId("59f43c5352e62394e2581f09") |

**Response**

Parameters for Success (Status code: 200)

| Parameter       | Type   | Description       |
| --------------- | ------ | ----------------- |
| itemId          | string | Unique item Id    |
| itemName        | float  | Name of the item  |
| itemPrice       | string | Price of the item |
| itemDescription | string | Ingredients used  |
| itemCalories    | float  | Calorie content   |

Parameters for Error (Status code: 404)

| Parameter | Type   | Description    |
| --------- | ------ | -------------- |
| message   | string | Item not found |

### 4. Get all items in the menu

**Request**

```
GET /menu/items
Content-Type: application/json
```

**Response**

Parameters for Success (Status code: 200) 

| Parameter       | Type   | Description       |
| --------------- | ------ | ----------------- |
| itemId          | string | Unique item Id    |
| itemName        | float  | Name of the item  |
| itemPrice       | string | Price of the item |
| itemDescription | string | Ingredients used  |
| itemCalories    | float  | Calorie content   |

Parameters for Error (Status code: 400)

| Parameter | Type   | Description   |
| --------- | ------ | ------------- |
| message   | string | Error Message |

### 5. Delete an item

**Request**

```
Delete /menu/item/:id
Content-Type: application/json
```

Parameters

| Parameter | Type   | Description    |
| --------- | ------ | -------------- |
| id        | string | Unique item Id |

**Response**

Parameters for Success (Status code: 200)

| Parameter | Type    | Description                                   |
| --------- | ------- | --------------------------------------------- |
| status    | boolean | Indicate that delete operation is successful. |
| message   | string  | Item deleted successfully                     |

Parameters for Error (Status code: 401)

| Parameter | Type   | Description                     |
| --------- | ------ | ------------------------------- |
| message   | string | Unauthorized to delete the item |

Parameters for Error (Status code: 404)

| Parameter | Type   | Description    |
| --------- | ------ | -------------- |
| message   | string | Item not found |