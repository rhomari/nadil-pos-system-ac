<template>
 
  <v-card
    class="py-8 px-6 text-center mx-auto ma-4 loginform"
    
    elevation="20"
    max-height="600"
    max-width="350"
    width="350"
   
    style="z-index:9999;margin-bottom:50px"
    
  >
  
      <user-select-component @user-selected="userSelected"></user-select-component>

   
    <v-sheet color="surface">
      <v-otp-input
        v-model="otp"
        type="password"
        variant="solo"
        length="6"  
       
        
      ></v-otp-input>
    </v-sheet>

    
    <div class="numpad mt-4" :disabled="userselected=='' ? '' : disabled">
      <div class="numpad-buttons">
        <button v-for="num in [1, 2, 3]" :key="'row1' + num" @click="appendNumber(num)" class="button">
          {{ num }}
        </button>
        <button v-for="num in [4, 5, 6]" :key="'row2' + num" @click="appendNumber(num)" class="button">
          {{ num }}
        </button>
        <button v-for="num in [7, 8, 9]" :key="'row3' + num" @click="appendNumber(num)" class="button">
          {{ num }}
        </button>
        <button class="button button-clear" @click="clearOtp">Vider</button>
        <button @click="appendNumber(0)" class="button">0</button>
      </div>
    </div>

    <v-btn
      class="my-4"
      color="purple"
      height="40"
      text="Verify"
      variant="flat"
      width="70%"
      @click="authenticate"
      :disabled="((otp.length < 6) || (userselected==null)) ? '' : disabled"
    >
      S'authentifier
    </v-btn>
    <p></p>
    <span v-if="autherror" style="color: darkred; font-size:12px">Erreur d'authenfication, merci de vérifier vos données de connexion.</span>

    <v-fab  size="32" color="red" icon="mdi-power" style="position: absolute;top: 96%;left: 45%;transform: translate(-50%, -50%);" @click="quitApp()" ></v-fab>
  </v-card>
  
</template>

<script>
import { getHostName, onUpdatefocused } from '@/Utils';
import UserSelectComponent from './UserSelectComponent.vue';
export default {
  components: { UserSelectComponent },
 
  data() {
    return {
      otp: "", 
      userselected:null,
      autherror:false,
      iswebview2:false,

    };
  },
 
  methods: {
    userSelected(user){
        this.userselected=user
        
        
    },
    quitApp(){
      if (window.nadilpos) {
        window.nadilpos.showMessage("close");
        return
      }   
      if (window.chrome.webview!=null)
        window.chrome.webview.postMessage("close")

    },
    appendNumber(num) {
      
      if (this.otp.length < 6) {
        this.otp += num.toString();
      }
    },
    clearOtp() {
      this.otp = "";
    },
    async authenticate() {
      
      if (this.userselected.Role === "SETUP"){
        this.$emit("auth-ok", this.userselected)
        return
      }
      this.autherror=false
      if (this.userselected==null){
        alert("Veuillez selectionner  un utilisateur");
        return
      }
      

      if (this.otp.length < 6) {
        alert("Veuillez enter votre mot de passe");
        return
      }
      let form = new FormData()
      form.append("user", this.userselected.Username)
      form.append("password", this.otp)
      await fetch("https://" + getHostName() + "/auth", {
        method:'post',
        body:form
      })
      .then(r=>{
        if (!r.ok) throw Error(r.statusText)
        return r.json()
        

      })
      .then(data=>{
        
        if (data.data == "OK"){
         
          this.$emit("auth-ok", this.userselected)
          return
        }
        this.autherror=true

      })
      .catch(ex=>{
        console.log(ex)
      })
      
    },
   
  },
};
</script>

<style scoped>
.numpad-buttons {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-top: 10px;
}


.button {
  padding: 10px;
  font-size: 16px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.button-clear {
  background-color: #FF5722;
  grid-column: span 2;
}

.button:active {
  background-color: #3e8e41;
}
@media only screen and (max-width: 574px) {
  .loginform{
    height: 525px;

    
    top:350px
  }
  
}
</style>
