<template>
  <div
    class="text-center pa-4"
    style="-webkit-user-select: none; -ms-user-select: none; user-select: none"
  >
    <v-dialog :model-value="modelValue" max-width="400" persistent>
      <v-card title="Editer un utilisateur">
        <template v-slot:actions>
          <div>
            <v-text-field
              label="Nom de l'utilisateur"
              required
              v-model="usertoedit.Username"
              disabled
              :rules="rules"
               @focus="onUpdatefocused"
            >
            </v-text-field>
          </div>
          <div>
            <v-select
              label="Role"
              required
              :items="['Gérant', 'Serveur']"
              v-model="role"
              :rules="rules"
            >
            </v-select>
          </div>
          <div>
            <v-text-field
              label="Mot de passe"
              required
              v-model="password"
              @keydown="isNumber"
              :type="showpassword ? 'text' : 'password'"
              clearable
              :append-icon="showpassword ? 'mdi-eye' : 'mdi-eye-off'"
              @click:append="showpassword = !showpassword"
              :rules="rules"
               @focus="onUpdatefocused"
            >
            </v-text-field>
          </div>
          <div>
            <v-text-field
              
              label="CODE"
              required
              v-model="usertoedit.Code"
              clearable
              type="password"
               @focus="onUpdatefocused"
            >
            </v-text-field>
          </div>

          <div>
            <v-btn @click="editUser()" icon="mdi-content-save"> </v-btn>

            <v-btn @click="closeDialog" icon="mdi-exit-to-app"> </v-btn>
          </div>
        </template>
      </v-card>
    </v-dialog>
  </div>
</template>
<script setup>
import { getHostName,  onUpdatefocused } from "@/Utils";
import { ref } from "vue";
// eslint-disable-next-line
const props = defineProps({ modelValue: Boolean, usertoedit: Object });
// eslint-disable-next-line
const emit = defineEmits(["update:modelValue", "alert", "inform"]);
const showpassword = ref(false);

const password = ref("");
const role = ref("");

function closeDialog() {
  emit("update:modelValue", false);
}
const rules = [
  (value) => {
    if (value) return true;

    return "ce champ est obligatoire";
  },
];

const isNumber = (event) => {
  if (
    password.value?.length === 6 &&
    event.key != "Backspace" &&
    event.key != "Delete"
  ) {
    event.preventDefault();
  }
  if (event.key.length === 1 && isNaN(Number(event.key))) {
    event.preventDefault();
  }
};
async function editUser() {
  if (
    props.usertoedit.Username == "" ||
    password.value == "" ||
    role.value == ""
  ) {
    alert("Tous les champs sont obligatoires");
    return;
  }
  let form = new FormData();
  form.append("username", props.usertoedit.Username);
  form.append("role", role.value);
  form.append("password", password.value);
  form.append("code", props.usertoedit.Code);
  form.append("id", props.usertoedit.ID);

  await fetch("https://" + getHostName() + "/edituser", {
    method: "post",
    body: form,
  })
    .then((r) => {
      if (!r.ok) {
        throw Error(r.statusText);
      }
      return r.json();
    })
    .then((data) => {
    
      
      emit("alert", data);
      emit("inform");
      
      closeDialog();
    })
    .catch((e) => {
      console.log(e);
    });
  props.usertoedit.Username = "" 
  password.value = "" 
  showpassword.value=false
  role.value = ""
  console.log(props.usertoedit.Username, role.value, showpassword.value,password.value)
}
</script>
<script>
export default {
  name: "NewCategoryDialogComponent",
  data() {},

  methods: {},
};
</script>

<style></style>
