function formatDate(dateString) {
    const seconds = Math.floor((new Date() - new Date(dateString)) / 1000);
    if (dateString.startsWith('0001-01-01')) {
        return '';
    }

    let interval = seconds / 31536000;
    if (interval > 1) return Math.floor(interval) + " years ago";
    
    interval = seconds / 2592000;
    if (interval > 1) return Math.floor(interval) + " months ago";
    
    interval = seconds / 86400;
    if (interval > 1) return Math.floor(interval) + " days ago";
    
    interval = seconds / 3600;
    if (interval > 1) return Math.floor(interval) + " hours ago";
    
    interval = seconds / 60;
    if (interval > 1) return Math.floor(interval) + " minutes ago";
    
    
    return Math.floor(seconds) + " seconds ago";
}

function displayDate(timestamp) {
    if (timestamp.startsWith('0001-01-01')) {
        return '';
    }
    const date = new Date(timestamp);
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    let hours = date.getHours();
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const ampm = hours >= 12 ? 'PM' : 'AM';
    hours = hours % 12;
    hours = hours ? hours : 12; // the hour '0' should be '12'
    return `${day}-${month}-${year} ${hours}:${minutes} ${ampm}`;
}


function ShortDate(dateString) {
    const seconds = Math.floor((new Date() - new Date(dateString)) / 1000);
    if (dateString.startsWith('0001-01-01')) {
        return '';
    }

    let interval = seconds / 31536000;
    if (interval > 1) return Math.floor(interval) + " years ago";
    
    interval = seconds / 2592000;
    if (interval > 1) return Math.floor(interval) + " months ago";
    
    interval = seconds / 86400;
    if (interval > 1) return Math.floor(interval) + " days ago";
    
    return "today";
}


function MessageLength(message){
    if (message.length > 18){
        return message.substring(0, 18) + '...';
    }
    return message;
}
function UserLength(user){
    if (user.length > 6){
        return user.substring(0, 6) + '.';
    }
    return user;
}
  
async function handlePostInteraction(event,distination) {
    event.preventDefault(); 
    const form = event.currentTarget;
    const postId = form.querySelector('input[name="postId"]').value;
    const isLike = form.classList.contains('like-form');
    
    try {
        const endpoint = isLike ? '/Data-PostLike' : '/Data-PostDisLike';
        const response = await fetch(`${endpoint}?postId=${postId}`, {
            method: 'POST',
            headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest'
            }
        });

        const data = await response.json();
        if (data.error && data.error === 'no userid') {
            return; // Do nothing if the error is "no userid"
        }
        const buttonsContainer = document.getElementById(`${distination}-post-${postId}`);
        if (!buttonsContainer) {
            console.error('Buttons container not found');
            return;
        }
        console.log(data);
        console.log(buttonsContainer);
        console.log(buttonsContainer.querySelector('.likes'));
        console.log(buttonsContainer.querySelector('.dislikes'));

        buttonsContainer.querySelector('.likes').textContent = data.LikeCount;
        buttonsContainer.querySelector('.dislikes').textContent = data.DislikeCount;
    } catch (error) {
        navigateToPage('Login');
    }   
}

