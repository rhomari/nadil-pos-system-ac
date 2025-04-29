<template>
  <div
    class="text-center pa-4"
    style="-webkit-user-select: none; -ms-user-select: none; user-select: none"
  >
    <v-dialog :model-value="modelValue" max-width="400" persistent>
      <v-card
        text="Indiquez le nom de la nouvelle catégorie"
        title="Nouvelle catégorie"
      >
        <template v-slot:actions>
          <div>
            <v-text-field
              label="Intitulé"
              required
              v-model="categorytoedit.categoryname"
               @focus="onUpdatefocused"
            >
            </v-text-field>
          </div>

          <div>
            <v-btn @click="saveCategory()" icon="mdi-content-save"> </v-btn>

            <v-btn @click="closeDialog" icon="mdi-exit-to-app"> </v-btn>
          </div>
        </template>
      </v-card>
    </v-dialog>
  </div>
</template>
<script setup>
import { getHostName, onUpdatefocused } from "@/Utils";

const props = defineProps({
  modelValue: Boolean,
  categorytoedit: { categoryname: "", categoryid: "" },
});
const emit = defineEmits(["update:modelValue"]);
function closeDialog() {
  props.categorytoedit.categoryname=''
  props.categorytoedit.categoryid=''
  emit("update:modelValue", false);
}

async function saveCategory() {
  const hostname = getHostName();
  const axios = require("axios").default;
  let requesturl = "https://" + hostname + "/saveNewCategory";
  let requestparams = { categoryname: props.categorytoedit.categoryname };
  if (props.categorytoedit.categoryid != "") {
    requesturl = "https://" + hostname + "/editCategory";
    requestparams = {
      categoryname: props.categorytoedit.categoryname,
      categoryid: props.categorytoedit.categoryid,
    };
  }

  await axios
    .post(requesturl, requestparams)
    .then((r) => {
      
      emit("alert", r.data);
    })
    .catch((e) => console.log(e));

  props.categorytoedit.categoryname = "";
  props.categorytoedit.categoryid = "";
  closeDialog();
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
