import { ref, nextTick, reactive} from "vue";
import Sortable from "sortablejs";
export function getHostName() {
  if (process.env.NODE_ENV !== "development")
    return window.location.hostname + ":" + window.location.port;
  return window.location.hostname + ":2304";
}

let sortableInstances = []; 


export const store = reactive({
  Member: null
})
export let ticketOK = false;
export const filtredMenu=ref([])
export  const initSortable = (id) => {
  
  nextTick(() => {
    
    sortableInstances.forEach((instance) => instance.destroy());
    sortableInstances = [];


    document.querySelectorAll(".sortable-grid").forEach((container, index) => {
      
      const sortable = Sortable.create(container, {
        animation: 150,
        ghostClass: "sortable-ghost",
        filter: ".card-element-no-drag",
        preventOnFilter: false,
        supportPointer: false,
        
        onEnd: (evt) => {
          
          const movedItem = filtredMenu.value.splice(evt.oldIndex, 1)[0];
          filtredMenu.value.splice(evt.newIndex, 0, movedItem);
          localStorage.setItem("filtredMenu" + id, JSON.stringify(filtredMenu.value));
          
         
        },
      }); 

      sortableInstances.push(sortable);
    });
  });
};
export const getSortedFiltredMenu = (id) => {
  const strlocalfiltredmenu = localStorage.getItem("filtredMenu" + id);
  
  
  const localfitredMenu = JSON.parse(strlocalfiltredmenu);
 
  const dbSet = new Set(filtredMenu.value.map(item => item.productid));
  
  let updatedLocalArray = localfitredMenu.filter(item => dbSet.has(item.productid));
  const localSet = new Set(updatedLocalArray.map(item => item.productid));
  const missingItems = filtredMenu.value.filter(item => !localSet.has(item.productid));

 
  updatedLocalArray = [...updatedLocalArray, ...missingItems];


  return updatedLocalArray;
}
export function onUpdatefocused() {
  if (window.chrome.webview!=null)
  window.chrome.webview.postMessage("vkeyboard")
}