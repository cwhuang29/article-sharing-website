# article-sharing-website
An article sharing website developed by Go.

## Overview
This project does not rely on any frontend framework, so it is a good entry point for backend engineers who want to build a whole website. With basic knowledge of JavaScript and CSS, You can start developing!
<br>
I chose [gin](https://github.com/gin-gonic/gin) as the backend web framework for its simplicity and high performance (it is also the most popular framework in Go, see [Top Go Web Frameworks](https://github.com/mingrammer/go-web-framework-stars).
<br>
For database ORM, I chose [gorm](https://github.com/go-gorm/gorm). It is a full-featured ORM with great community support and [easy to follow documentation](https://gorm.io/docs/).
Besides, if you chose `sqlite` as the database driver, then you can get rid of the database burden (the data will be stored in a file `tmp.db` in the root of project) and focus on the backend development part.

## Setup
As the register feature is not yet implemented, there is a fake user registered with email `admin@gmail.com` and password `a1234567`. This will be removed in the future.
<br>
For the same reason, the upload article function is hidden from the website interface. Nevertheless, you can open the webpage via `http://127.0.0.1/admin/create/article`.
<br>

#### Local
If you want to use `MySQL` as the database driver, you have to create a database in advance and set the corresponding database name, username, and password in `config.yml`.
```
go build -o web cmd/main.go
APP_PORT=8080 DB_HOST=127.0.0.1 ./web
```

#### Docker
```
docker run -d \
    --name db \
    -e MYSQL_DATABASE=inews \
    -e MYSQL_ROOT_PASSWORD=a1234567 \
    -e MYSQL_USER=user01 \
    -e MYSQL_PASSWORD=a1234567 \
    mysql:5.7.32 \
    mysqld --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

docker run -d \
    --name web \
    --link db:db \
    -e APP_PORT=80\
    -e DB_HOST=db \
    -p 80:80 \
    cwhuang29/article-sharing-website:v1
```

## TODO
- [ ] Register
- [ ] Support Google login
- [ ] Logout
- [ ] Search articles based on tags (currently only the category tags have this feature)
- [ ] Pagination (in the article overview page)
- [ ] Logger
- [ ] Likes/Dislikes buttons
- [ ] Preview feature (before submitting the article)
- [ ] Security issues about uploading files (e.g. check files type andfiles size)

## Demo
Here is a live demo: [inews](http://18.179.7.226/) (hosting on AWS)
### Feature - article list
![Articles List](demo/articles-list.png)
### Feature - browse
![Browse](demo/browse.png)
### Feature - edit
![Edit](demo/edit.png)

