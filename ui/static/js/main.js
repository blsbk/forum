
var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

document.addEventListener("DOMContentLoaded", function () {
    
    let isLiked = document.getElementById("isLiked").innerHTML;
	let likeCnt = +document.getElementById("likeCount").textContent;
	const likeButton = document.getElementById("likeButton");
	const likeIcon = document.getElementById("likeIcon");
	const likeCountElement = document.getElementById("likeCount");
	isLiked = (isLiked?.toLowerCase?.() === 'true');
    
    let postID = +document.getElementById("postID").innerHTML;

	likeButton.addEventListener("click",(event) => {

        // Update the like icon based on the liked state
        if (!isLiked) {
            likeIcon.src = "/static/img/like.png";
            likeCnt++;
        } else if (isLiked) {
            likeIcon.src = "/static/img/unlike.png";
            likeCnt--;
        }
        isLiked = !isLiked


        // Update the like count in the HTML
        likeCountElement.textContent = likeCnt;
        let body = { likeCount: likeCnt, postID: postID, isLiked: isLiked }
        let url = "/post/like"
        submitLike(body, url)
       
    })

});

function submitLike(body, url) {
    fetch(url, {
        method: "POST",
        body: JSON.stringify(body),
        headers: {
            "Content-Type": "application/json"
        }
    })
    .then(response => {
        if (response.ok) {
            location.reload()
            console.log("Succesfully liked")
        } else {
            console.error("Failed to like");
        }
    })
    .catch(error => {
        console.error(error);
    });
}

document.addEventListener("DOMContentLoaded", function () {
    
   
	const commentLikeButton = document.querySelectorAll("#commentLikeButton")
	
    
    
    let postID = +document.getElementById("postID").innerHTML;

    commentLikeButton.forEach((button) => {
        let isCommentLiked = document.getElementById("isCommentLiked").innerHTML;
        let commentLikeCnt = +document.getElementById("commentLikeCount").textContent;
        const commentLikeIcon = document.getElementById("commentLikeIcon");
	    const commentLikeCountElement = document.getElementById("commentLikeCount");
	    isCommentLiked = (isCommentLiked?.toLowerCase?.() === 'true');
        let commentID = +document.getElementById("commentID").innerHTML;
        
        button.addEventListener("click", (event) => {
            // Update the like icon based on the liked state
            if (!isCommentLiked) {
               commentLikeIcon.src = "/static/img/thumbUpClicked.png";
               commentLikeCnt++;
           } else if (isCommentLiked) {
               commentLikeIcon.src = "/static/img/thumbUpUnclicked.png";
               commentLikeCnt--;
           }
           isCommentLiked = !isCommentLiked
   
   
           // Update the like count in the HTML
           commentLikeCountElement.textContent =commentLikeCnt;
           let body = { commentLikeCount: commentLikeCnt, commentID: commentID, postID: postID, isCommentLiked: isCommentLiked }
           let url = "/post/commentLike"
           submitCommentLike(body, url)
       })
    })

});

function submitCommentLike(body, url) {
    fetch(url, {
        method: "POST",
        body: JSON.stringify(body),
        headers: {
            "Content-Type": "application/json"
        }
    })
    .then(response => {
        if (response.ok) {
            location.reload()
            console.log("Succesfully liked comment")
        } else {
            console.error("Failed to like comment");
        }
    })
    .catch(error => {
        console.error(error);
    });
}
