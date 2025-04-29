<template>
    <div class="user-select">
        <div class="dropdown" >
        <button class="dropdown-button" @click="toggleDropdown">
          <span v-if="selectedUser">
            
            
                <b style="font-weight: bold;">{{selectedUser.Username.toUpperCase()}}</b>
                
            
            
          </span>
          <span v-else>Sélectionnez un utilisateur</span>
          <svg class="chevron" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 10.94l3.72-3.72a.75.75 0 111.06 1.06l-4.25 4.25a.75.75 0 01-1.06 0L5.23 8.27a.75.75 0 01-.02-1.06z" clip-rule="evenodd" />
          </svg>
        </button>
        <transition name="fade">
          <div v-if="showDropdown" class="dropdown-list">
            <div
              v-for="user in users"
              :key="user.ID"
              class="dropdown-item"
              @click="selectUser(user)"
            >
             
             
              {{ user.Username.toUpperCase() }}
            </div>
          </div>
        </transition>
      </div>
    </div>
  </template>
  
  <script>
  import { ref,defineEmits, onMounted } from 'vue';
  import { getHostName } from '@/Utils';
  
  export default {
    name: 'UserSelect',
    
    setup(props,{emit}) {
      
      
      const users = ref(null)
      const selectedUser = ref(null);
      const showDropdown = ref(false);
      
      const toggleDropdown = () => {
        showDropdown.value = !showDropdown.value;
      };
      onMounted(async () => {
        await getUsers();
        
      }); 
     
      const selectUser = (user) => {
        emit("user-selected", user)
        selectedUser.value = user;
        
       
        showDropdown.value = false;
        
        
      };
      async function getUsers() {
        
        
        await fetch("https://" + getHostName() + "/getusers")
        .then((r) => {
            if (!r.ok) throw Error(r.statusText);
                return r.text();
        })
        .then((data) => {
            users.value = JSON.parse(data);

            if (users.value== null || !users.value.some(obj => obj.Role == "Gérant" )){
              users.value= [{ID: '-1', Username: 'SETUP', Role: 'SETUP', Code: ''}]
            }
            
      
         
        })
        .catch((e) => console.log(e));
    }
      const getInitials = (name) => {
        
        
       
      };
  
      return {
        users,
        selectedUser,
        showDropdown,
        toggleDropdown,
        selectUser,
        getInitials,
      };
    },
  };
  </script>
  
  <style scoped>
  .user-select {
    font-family: Arial, sans-serif;
    max-width: 500px;
    margin: 10px;
  }
  
  .dropdown {
    position: relative;
    display: inline-block;
    width: 100%;
  }
  
  .dropdown-button {
    width: 100%;
    padding: 12px;
    background: #f9f9f9;
    border: 1px solid #ddd;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    cursor: pointer;
    font-size: 16px;
    color: #333;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    transition: background 0.3s ease;
  }
  
  .dropdown-button:hover {
    background: #f0f0f0;
  }
  
  .dropdown-button .chevron {
    width: 20px;
    height: 20px;
  }
  
  .dropdown-list {
    position: absolute;
    top: 110%;
    left: 0;
    width: 100%;
    background: #fff;
    border: 1px solid #ddd;
    border-radius: 8px;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
    overflow: hidden;
    max-height: 250px;
    overflow-y: auto;
    z-index: 1000;
  }
  
  .dropdown-item {
    padding: 10px 15px;
    cursor: pointer;
    display: flex;
    align-items: center;
    transition: background 0.2s;
    color: black !important;
  }
  .dropdown-item:hover {
    color: tomato !important;
  }
  
  .user-image, .user-initials {
    width: 35px;
    height: 35px;
    border-radius: 50%;
    margin-right: 10px;
  }
  
  .user-image {
    object-fit: cover;
  }
  
  .user-initials {
    background-color: #666;
    color: #fff;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 16px;
  }
  
  .fade-enter-active, .fade-leave-active {
    transition: opacity 0.3s ease;
  }
  .fade-enter-from, .fade-leave-to {
    opacity: 0;
  }
  </style>
  