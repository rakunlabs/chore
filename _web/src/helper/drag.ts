const moveElement = (move:HTMLElement, X: boolean, Y: boolean, fn: (X: number, Y:number)=>void, down: ()=>void = ()=>void{} ) => {
  let posX = 0;
  let posY = 0;
  let posChangeX = 0;
  let posChangeY = 0;

  move.onmousedown = dragMouseDown;

  function dragMouseDown(e:MouseEvent) {
    if (e.target != move) {
      return;
    }
    e.preventDefault();
    // get the mouse cursor position at startup:
    if (Y) {
      posY = e.clientY;
    }
    if (X) {
      posX = e.clientX;
    }
    document.onmouseup = closeDragElement;
    // call a function whenever the cursor moves:
    move.onmousemove = elementDrag;
    down();
  }

  function elementDrag(e:MouseEvent) {
    // console.log(e.currentTarget);
    if (e.target != move) {
      return;
    }

    e.preventDefault();
    // calculate the new cursor position:
    if (Y) {
      posChangeY = posY - e.clientY;
      posY = e.clientY;
    }
    if (X) {
      posChangeX = posX - e.clientX;
      posX = e.clientX;
    }
    // set the element's new position:
    fn(posChangeX, posChangeY);
  }

  function closeDragElement() {
    // stop moving when mouse button is released:
    document.onmouseup = null;
    move.onmousemove = null;
  }
};

export { moveElement };
