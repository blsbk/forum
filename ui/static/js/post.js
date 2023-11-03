
let likeCnt = +document.getElementById("likeCount").textContent;
let dislikeCnt = +document.getElementById("dislikeCount").textContent;


const commentLikeButtons = document.querySelectorAll(".commentLikeButton")
const commentDislikeButtons = document.querySelectorAll(".commentDislikeButton")
const commentLikeCntElement = document.querySelectorAll(".commentLikeCount");
const commentDislikeCntElement = document.querySelectorAll(".commentDislikeCount");
const commentLikeCnt = Array.from(commentLikeCntElement).map(element => +element.innerText);
const commentDislikeCnt = Array.from(commentDislikeCntElement).map(element => +element.innerText);

const postID = +document.getElementById("postID").innerHTML;


document.addEventListener("DOMContentLoaded", function () {
    
    let isLiked = document.getElementById("isLiked").innerHTML;
	const likeButton = document.getElementById("likeButton");
	isLiked = (isLiked?.toLowerCase?.() === 'true');
    
	likeButton.addEventListener("click",(event) => {

        if (!isLiked) {
            likeCnt++;
        } else if (isLiked) {
            likeCnt--;
        }
        isLiked = !isLiked

        let body = { likeCount: likeCnt, postID: postID, isLiked: isLiked }
        let url = "/post/like"
        submitLike(body, url, "like")
    })
});

document.addEventListener("DOMContentLoaded", function () {
    
    let isDisliked = document.getElementById("isDisliked").innerHTML;
	const dislikeButton = document.getElementById("dislikeButton");
	isDisliked = (isDisliked?.toLowerCase?.() === 'true');
    
	dislikeButton.addEventListener("click",(event) => {

        if (!isDisliked) {
            dislikeCnt++;
        } else if (isDisliked) {
            dislikeCnt--;
        }
        isDisliked = !isDisliked

        let body = { dislikeCount: dislikeCnt, postID: postID, isDisliked: isDisliked }
        let url = "/post/dislike"
        submitLike(body, url, "dislike")
    })
});

function submitLike(body, url, msg) {
    fetch(url, {
        method: "POST",
        body: JSON.stringify(body),
        headers: {
            "Content-Type": "application/json"
        }
    })
    .then(response => {
        if (response.ok) {
            location.href = location.href;
            console.log("Succesfull post "+msg)
        } else {
            console.error("Failed post "+msg);
        }
    })
    .catch(error => {
        console.error(error);
    });
}

commentLikeButtons.forEach((button, index) => {
    let commentID = +button.getAttribute("comment-id");
    let isCommentLiked = button.getAttribute("comment-liked");
    isCommentLiked = (isCommentLiked?.toLowerCase?.() === 'true');
    
    button.addEventListener("click", () => {
        if (!isCommentLiked) {
            commentLikeCnt[index]++;
        } else if (isCommentLiked) {
            commentLikeCnt[index]--;
        }
        isCommentLiked = !isCommentLiked;
        let body = {commentLikeCount: commentLikeCnt[index], commentID: commentID, postID: postID, isCommentLiked: isCommentLiked}
        let url = "/post/commentLike"
        submitCommentLike(body, url, "like");
    })
})

commentDislikeButtons.forEach((button, index) => {
    let commentID = +button.getAttribute("comment-id");
    let isCommentDisliked = button.getAttribute("comment-disliked");
    isCommentDisliked = (isCommentDisliked?.toLowerCase?.() === 'true');
    
    button.addEventListener("click", () => {
        if (!isCommentDisliked) {
            commentDislikeCnt[index]++;
        } else if (isCommentDisliked) {
            commentDislikeCnt[index]--;
        }
        isCommentDisliked = !isCommentDisliked;

        let body = {commentDislikeCount: commentDislikeCnt[index], commentID: commentID, postID: postID, isCommentDisliked: isCommentDisliked}
        let url = "/post/commentDislike"
        submitCommentLike(body, url, "dislike");
    })
})

function submitCommentLike(body, url, msg) {
    fetch(url, {
        method: "POST",
        body: JSON.stringify(body),
        headers: {
            "Content-Type": "application/json"
        }
    })
    .then(response => {
        if (response.ok) {
            location.href = location.href;
            console.log("Succesfull comment "+msg);
        } else {
            console.error("Failed comment "+msg);
        }
    })
    .catch(error => {
        console.error(error);
    });
}
