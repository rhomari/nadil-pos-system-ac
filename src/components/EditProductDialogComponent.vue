<template>
  <div class="text-center pa-4">
    <v-dialog :model-value="modelValue" max-width="400" persistent>
      <v-card text="Modification du produit" title="Modification">
        <template v-slot:actions>
          <div>
            <v-text-field
              :label="getCategory(editproduct.category)"
              prepend-icon="mdi-basket"
              disabled
              density="compact"
               @focus="onUpdatefocused"
            >
            </v-text-field>
            <v-text-field
              label="Intitulé"
              v-model="editproduct.text"
              prepend-icon="mdi-food"
              density="compact"
              :rules="rules"
               @focus="onUpdatefocused"
            >
            </v-text-field>
            <v-text-field
              label="Prix"
              v-model="editproduct.price"
              prepend-icon="mdi-cash-multiple"
              type="number"
              density="compact"
              :rules="rules"
               @focus="onUpdatefocused"
            ></v-text-field>
            <v-file-input
              accept="image/png, image/jpeg, image/bmp"
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
            <v-checkbox
              v-model="editproduct.isfavorite"
              color="indigo"
              label="Favori"
            ></v-checkbox>
          </div>

          <div>
            <v-btn v-on:click="editProduct" icon="mdi-content-save"> </v-btn>

            <v-btn v-on:click="closeDialog" icon="mdi-exit-to-app"> </v-btn>
          </div>
        </template>
      </v-card>
    </v-dialog>
  </div>
</template>
<script setup>
import { ref, watch } from "vue";

import { getHostName } from "@/Utils";
const props = defineProps({
  modelValue: Boolean,
  editproduct: Object,
  categories: Object,
});
const emit = defineEmits(["update:modelValue"]);
let categoryname;

const imagefile = ref("");
let preview = ref("");

async function editProduct() {
  if (props.editproduct.price<0) {
    alert("Le prix ne peut pas prendre une valeur négative.")
    return
  }

  const hostname = getHostName();

  if (props.editproduct.text == "" || props.editproduct.price == "") {
    alert("Veuillez remplir les champs obligatoires");
    return;
  }

  const form = new FormData();
  form.append("imagefile", imagefile.value);
  form.append(
    "details",
    JSON.stringify({
      productid: props.editproduct.productid,
      isfavorite: props.editproduct.isfavorite,
      category: props.editproduct.category,
      text: props.editproduct.text,
      price: props.editproduct.price,
    })
  );
  await fetch("https://" + hostname + "/editProduct", {
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
function closeDialog() {
  imagefile.value=null
  preview.value=null
  emit("update:modelValue", false);
}
function getCategory(id) {
  props.categories.forEach((c) => {
    if (c.categoryid == id) {
      categoryname = c.categoryname;
      return categoryname;
    }
  });

  return categoryname;
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
};
</script>

<style></style>
