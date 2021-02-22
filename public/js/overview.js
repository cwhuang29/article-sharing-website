const fetchArticle = async (url = '', data = {}) => {
    const response = await fetch(url, {
        method: 'GET', // *GET, POST, PUT, DELETE, etc.
        mode: 'cors', // no-cors, *cors, same-origin
        mode: 'same-origin',
        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        credentials: 'same-origin', // include, *same-origin, omit
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8'
        },
        redirect: 'follow', // manual, *follow, error
        referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
    });
    return response;
};

const DECREASE_IMAGE_HEIGHT = 950;
const imageDisplayToggle = () => {
	let imgDisplay = true;
	if (window.innerWidth < DECREASE_IMAGE_HEIGHT) {
		imgDisplay = false;
		[...document.querySelectorAll('.article-list-img-h')].forEach(img => {
            img.style.height = 132;
		});
	}
	window.addEventListener('resize', _ => {
		if (imgDisplay == true && window.innerWidth < DECREASE_IMAGE_HEIGHT) {
			imgDisplay = false;
			[...document.querySelectorAll('.article-list-img-h')].forEach(img => {
                img.style.height = 150;
			});
		} else if (imgDisplay == false && window.innerWidth >= DECREASE_IMAGE_HEIGHT) {
			imgDisplay = true;
			[...document.querySelectorAll('.article-list-img-h')].forEach(img => {
                img.style.height = 172;
			});
		}
	});
};

onDOMContentLoaded = (function(){
    /*
    const checkStatus = async (resp) => {
        if (resp.status >= 400) {
            return Promise.reject(resp);
        } else {
            return Promise.resolve(resp);
        }
    };
    const showArticle = async (resp) => {
        c(resp.redirected, resp.url);
        if (response.redirected) { // Should set the Location header in server side
            window.location.href = response.url;
        }
        // resp.redirect(); // ???????????
        // resp.json().then(function(data) {
        //     c(data);
        // })
        return Promise.resolve();
    };
    const handleErr = (resp) => {
        if (resp.status == 400 || resp.status == 404) {
            resp.json().then((data) => {
                showErrMsg(data);
            })
        } else {
            showErrMsg("<div><p><strong> Some severe errors occurred !</strong></p><p>Please reload the page and try again.</p></div>");
        }
    };
    const getArticleFullContent = (id) => {
        document.body.style.cursor = 'wait!important';
        fetchArticle('/articles/browse?' + new URLSearchParams({"id": id})) // /articles/browse?id=<id>
            .then(checkStatus)
            .then(showArticle)
            .catch(handleErr)
            .finally(_ => {
                document.body.style.cursor = 'default';
            });
    };
    */

    [...document.querySelectorAll('.tile > .title, .tile > .subtitle, .tile > .columns')].forEach(tile => {
        tile.addEventListener('click', (e) => {
            // getArticleFullContent(tile.parentNode.children[0].dataset.articleid);
            window.location.href = '/articles/browse?articleId=' + tile.parentNode.children[0].dataset.articleid;
        });
    });

	imageDisplayToggle();
})();
