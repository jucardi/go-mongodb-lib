# Pages bundle

This bundle goes along with the `github.com/jucardi/go-mongo-lib/mgo/query_pages.go` extension.

In a very generic way, this bundle allows the use of page parameters as query strings in HTTP request, and handles the logic necessary to retrieve the proper items requested by a page from a MongoDb instance.

### Features this bundle provides:

- **Extension for `github.com/gin-gonic/gin`**<br>
  Automatically create a `*pages.Page` struct with page request information extracted directly from a query string added to an HTTP request, obtained from a `*gin.Context`.\
  For more information, see `CreateFromContext(c *gin.Context, defaultPage ...*Page) *Page` in package `github.com/jucardi/go-mongo-lib/pages`

- **Extension for `gopkg.in/mgo.v2`** *(through `github.com/jucardi/go-mongo-lib/mgo`)*<br> 
  Adds helper functions to `mgo.IQuery` to easily obtain a page result from a collection by using the provided page information before retrieving the final results.\
  For more information, see the following `IQuery` extension functions in `github.com/jucardi/go-mongo-lib/mgo/components/mgo`:
  - `Page(page ...*pages.Page) IQuery`
  - `WrapPage(result interface{}, page ...*pages.Page) (*pages.Paginated, error)`

### Usage

The query strings used are the following:
- *`'page'`*: Indicates the page number to be requested.
- *`'size'`*: Indicates the page size (amount of items per page).
- *`'sort'`*: *(optional)* Multiple `sort` values may be passed. Indicates the fields to use to sort the sample before retrieving a page. Must match a key in the mongo document. Use '-' at the beginning for reverse order. Eg "-name".

**Example:**
```
curl http://user-service:1234/users?page=5&size=10&sort=firstname&sort=-lastname
```

The query above will send a request to the route `http://some-service:1234/path-to-route`, which internally add the necessary information to the request for MongoDb to first, sort by `firstname` then by `lastname` in reverse order, and return items from index 40 to index 49 (page 5, size 10)

The result should look like this

```json
{
    "count": 10,          // Indicates the amount of items retrieved. Normally equal to page size unless page >= total pages
    "total_pages": 35,    // Total pages in the full result set
    "total_count": 342,   // Total items in the full result set.
    "size": 10,           // The requested page size.
    "current_page": 5,    // The requested page number.
    "content": [ . . . ]  // Array of JSON documents that belong to the requested page.
}
``` 
<br>

### Using pagination in a `mgo` repository implementation.

If already using the `github.com/jucardi/go-mongo-lib/mgo` wrapper, continue to the next section, otherwise follow the steps below.

1) Replace the import for `gopkg.in/mgo.v2` with `github.com/jucardi/go-mongo-lib/mgo`
2) Add the import for the mongo middleware `github.com/jucardi/go-mongo-lib/middleware/mongo`

The mongo middleware opens a connection to the database o startup, and this connection remains open to avoid having to do the dial for every single request. The session is automatically cloned within the request context
so each concurrent call has their own Session to work with. The function `mongo.GetDb()` returns a clone of the session and the database. This facilitates:

```Go
    session, db := getDb()
    defer session.Close()
```

The new `ISession` and `IDatabase` interfaces implement the same functions found in `*mgo.Database` and `*mgo.Session`, so it should be a seamless change.

**Note**: If using direct references to the functional structs in mgo such as `*mgo.Database`, `*mgo.Session`, `*mgo.Collection`, do the following replaces in the code:
  - `*mgo.Database` with `mgo.IDatabase`
  - `*mgo.Collection` with `mgo.ICollection`
  - `*mgo.Query` with `mgo.IQuery`
  - `*mgo.Session` with `mgo.ISession`
  - `*mgo.Database` with `mgo.IDatabase`
  - `*mgo.Bulk` with `mgo.IBulk`
  - `*mgo.Iter` with `mgo.IIter`
This will make reference to the new Interface wrappers in `github.com/jucardi/go-mongo-lib/mgo` created for `gopkg.in/mgo.v2`*

<br>

#### Adding a function to the repository which queries to MongoDb and wraps the results in `*pages.Paginated`

To achieve this, after doing any queries to Mongo, the `mgo.IQuery` implements a function `WrapPage`, which receives a pointer to the array where the results will be stored, and a variadic `...*pages.Page` arg (making the page argument option).
This will automatically:
- Calculate the size the sample that matches the query at the state before calling `WrapPage`
- Calculate the total amount of pages in the sample with the provided *page size*
- `Skip` the first N records (N obtained by multiplying the provided *page size* and *page number*)
- `Limit` the results to the provided *page size*
- Wrap the provided array pointer in a a new instance of `*pages.Paginated` and append the Paginated information (page size, page number, total items retrieved, total items in the query, total pages)

**Example**
```Go
func (r *repository) GetAll(page ...*pages.Page) (*pages.Paginated, error) {
    session, db := getDb()
    defer session.Close()
    var result []*User
    return db.C("users").Find(bson.M{}).WrapPage(&result, page...)
}
```
> *Any filtering, sorting or any other query operation that returns a query can be used before invoking `WrapPage`*

<br>
<br>

### Adding the pages `github.com/gin-gonic/gin` bundle to a git route handler.

Simply create the page object by doing `page := pages.CreateFromContext(c)`. This will create the `*pages.Page` from the query strings.

**Example**
```Go
func getUsers(c *gin.Context) {
    page := pages.CreateFromContext(c)
    if ret, err := users.Repo().GetAll(page); err != nil {
        c.IndentedJSON(err.Code, err)
    } else {
        c.IndentedJSON(http.StatusOK, ret)
    }
}
```
<br>
<br>

### Creating a friendly Response Struct for Golang RestClients implementations to consume the `*pages.Paginated` result

*In most cases, this step is not necessary, like creating an API that will be consumed by a client written in a different technology, such as a React application. This is only recommended when creating RestClient implementation in Golang*

The `pages.Paginated` struct embeds `pages.PaginatedBase` which defines the basic fields of the Paginated object (all but the `json:"content"` field). This allows to easily create any implementation
With the proper array type in a very simple way.

To do this simply declare a struct that will have embedded a `pages.PaginatedBase` ("inherit" from it), and add an array field of the object type. The field name can be anything, as long as the json mapping for that field is `content`

**Example**
```Go
type PaginatedUsers struct {
    pages.PaginatedBase
    Items []*Users `json:"content"`
}
```
> After doing any `http` operation that yields a response (using a `*http.Response` type for this example) The paginated object can be easily deserialized
```Go
    var resp *http.Response
    
    . . . // Rest Client logic here

    paginated := &PaginatedUsers{}
    respBytes, _ := ioutil.ReadAll(resp.Body)
    json.Unmarshal(respBytes, paginated)
```

