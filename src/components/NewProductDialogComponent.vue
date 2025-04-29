<template>
  <div class="text-center pa-4">
    <v-dialog :model-value="modelValue" max-width="400" persistent>
      <v-card text="Indiquez le nom du nouveau produit" title="Nouveau produit">
        <template v-slot:actions>
          <div>
            <v-text-field
              :label="parentcategory.categoryname"
              prepend-icon="mdi-basket"
              disabled
              density="compact"
              
            >
            </v-text-field>
            <v-text-field
              label="Intitulé"
              v-model="newproductname"
              prepend-icon="mdi-food"
              density="compact"
              :rules="rules"
               @focus="onUpdatefocused"
            >
            </v-text-field>
            <v-text-field
              label="Prix"
              v-model="price"
              prepend-icon="mdi-cash-multiple"
              type="number"
              density="compact"
              :rules="rules"
               @focus="onUpdatefocused"
            ></v-text-field>
            <v-file-input
              accept="image/png, image/jpeg, image/bmp, image/webp"
              label="image"
              @change="makePreview"
              density="compact"
              :rules="rules"
              v-model="imagefile"
              prepend-icon="mdi-camera"
            >
            </v-file-input>
            <v-img
              :width="100"
              :height="100"
              style="float: right; margin-bottom: 10px; border-radius: 5px"
              cover
              :src="preview"
            ></v-img>
            <v-switch
              density="compact"
              color="green"
              v-model="isFavorite"
              label="Favori"
            ></v-switch>
          </div>

          <div>
            <v-btn @click="saveNewProduct()" icon="mdi-content-save"> </v-btn>

            <v-btn @click="closeDialog" icon="mdi-exit-to-app"> </v-btn>
          </div>
        </template>
      </v-card>
    </v-dialog>
  </div>
</template>
<script setup>
import { getHostName, onUpdatefocused } from "@/Utils";
import { ref } from "vue";

const props = defineProps({
  modelValue: Boolean,
  parentcategory: { categoryname: String, categoryid: Number },
});
const emit = defineEmits(["update:modelValue"]);
const isFavorite = ref(false);
const newproductname = ref("");
const price = ref(0);
const preview = ref(null)
const imagefile = ref("");
function closeDialog() {
  
  newproductname.value=""
  price.value=''
  imagefile.value=null
  preview.value=null

  emit("update:modelValue", false);
}

async function saveNewProduct() {
  if (price.value<0) {
    alert("Le prix ne peut pas prendre une valeur négative.")
    return
  }
  const hostname = getHostName();
  if (newproductname.value == "" || imagefile.value == "") {
    alert("Veuillez remplir les champs obligatoires");
    return;
  }
  console.log("sending...");

  const form = new FormData();
  form.append("imagefile", imagefile.value);
  form.append(
    "details",
    JSON.stringify({
      isfavorite: isFavorite.value,
      category: props.parentcategory.categoryid,
      text: newproductname.value,
      price: price.value,
    })
  );
  await fetch("https://" + hostname + "/addNewProduct", {
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

      closeDialog();
    })
    .catch((e) => {
      console.log(e);
    });
}

function setPreview(file) {
  const r = new FileReader();
  r.onload = (e) => {
    preview.value = e.target.result;
  };
  r.readAsDataURL(file);
}
function makePreview(e) {
  if (!e.target.files[0]) {
    return;
  }
  
  setPreview(e.target.files[0]);
}
</script>

<script>
export default {
  data() {
    return {
     

      rules: [
        (value) => {
          if (value) return true;

          return "ce champ est obligatoire";
        },
      ],
    };
  },
  methods: {
    
  },
};
</script>

<style></style>
