module qilin-api

require (
	cloud.google.com/go v0.34.0 // indirect
	github.com/ProtocolONE/authone-jwt-verifier-golang v0.0.0-20190329122021-aa7178c82afb
	github.com/ProtocolONE/rbac v0.0.0-20190327162713-7c9ef6330a8a
	github.com/casbin/redis-adapter v0.0.0-20190105032110-b36d844dade5
	github.com/centrifugal/centrifuge-go v0.2.3
	github.com/centrifugal/gocent v2.0.2+incompatible
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20181014144952-4e0d7dc8888f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/google/go-cmp v0.2.0 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/jinzhu/gorm v1.9.2
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/jinzhu/now v0.0.0-20181116074157-8ec929ed50c3 // indirect
	github.com/jpillora/backoff v0.0.0-20180909062703-3050d21c67d7 // indirect
	github.com/karlseguin/ccache v2.0.2+incompatible // indirect
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/labstack/echo v3.3.5+incompatible
	github.com/labstack/echo/v4 v4.0.0
	github.com/labstack/gommon v0.2.8
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/lib/pq v1.0.0
	github.com/lunny/html2md v0.0.0-20181018071239-7d234de44546
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/microcosm-cc/bluemonday v1.0.2
	github.com/mitchellh/mapstructure v1.1.2
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/valyala/fasttemplate v1.0.1 // indirect
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1
	golang.org/x/crypto v0.0.0-20190325154230-a5d413f7728c // indirect
	golang.org/x/sys v0.0.0-20190322080309-f49334f85ddc // indirect
	google.golang.org/appengine v1.4.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/go-playground/validator.v9 v9.26.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1
