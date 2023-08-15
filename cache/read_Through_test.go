package cache

//func TestReadThroughCacheV1_Get(t *testing.T) {
//	var c1 ReadThroughCache = ReadThroughCache{
//		LoadFunc: func(ctx context.Context, key string) (any, error) {
//			if strings.HasSuffix(key, "order_1") {
//				//加载 order
//			} else if strings.HasSuffix(key, "user") {
//				//加载 user
//			}
//		},
//	}
//	var c2 ReadThroughCacheV1[User]
//
//	val, _ := c1.Get(context.Background(), "user_v1")
//	u := val.(User)
//
//	x, _ := c2.Get(context.Background(), "user_v1")
//	t.Log(x.Name)
//
//}
//
//type User struct {
//	Name string
//}
