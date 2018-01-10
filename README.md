# twitch-plays-factorio

This is the central Repo for all things "twitch-plays-factorio"

## Useful scripts

### Easy mouse buttons
Paste this into the console & press enter.
```javascript
$('#root-player').hide();
$(".player-overlay").on('mousedown',function(e) {
  e.preventDefault();
  $('[data-a-target=chat-input]').val( (e.button == 0?'p':'s')+"("+parseInt(e.offsetX / e.target.offsetWidth*1920)+","+parseInt(e.offsetY/ e.target.offsetHeight*1080)+"").focus();
});
```
You can use your mouse button to send clicks! Thanks [@fumai](https://github.com/fuami)!
