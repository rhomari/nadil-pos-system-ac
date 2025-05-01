/* eslint-disable */
<template>
  
 

  <v-responsive class="border rounded no-select">
    <div class="text-center">
   

    <v-overlay  v-if="authenticated==false" v-model="overlay" persistent :style="getOverlayStyle()">
    
    </v-overlay>
    <auth-num-pad-component v-if="authenticated==false"  class="centerstuff"  @auth-ok="authOK"></auth-num-pad-component>
    
  </div>
    <v-app :theme="theme">
      <v-app-bar :elevation="4" density="compact">
        <v-menu>
      <template v-slot:activator="{ props }">
        <v-btn
          style="font-size:10px"
          color="primary"
          v-bind="props"
        >
          Thème
        </v-btn>
      </template>
      <v-list>
        <v-list-item
          v-for="(t, index) in themes"
          :key="index"
          style="cursor:default"
          @click="changeTheme(t.title)"
        >
          <v-list-item-title style="font-size:12px; font-weight:bold">{{ t.title }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
        <template v-slot:prepend>
          <v-app-bar-nav-icon @click="showMenu"></v-app-bar-nav-icon>
        </template>

        <v-app-bar-title>Nadil-POS AC</v-app-bar-title>
        <template v-if="authenticated" v-slot:append>
          <span ref="authuser">{{loggeduser.Username.toUpperCase()}}</span>
          
          <v-btn  @click="authOff();logout()"> <v-icon size="18px">mdi-logout</v-icon></v-btn>
          <v-btn value="quit" @click="minimizeApp()" style="width: 48px;" class="NOMOBILE">
          <v-icon>mdi-window-minimize</v-icon>

          
        </v-btn>
          <v-btn value="quit" color="red" @click="quitApp()" style="width: 48px;" >
          <v-icon size="24px">mdi-power</v-icon>

          
        </v-btn>
        </template>
      </v-app-bar>

      <v-main>
        <v-container>
    
          <v-card
            class="menu"
            style="float: left; overflow: auto; max-height: 80vh; height: 80vh"
          >
            <v-tabs
              v-model="tab"
              align-tabs="end"
              next-icon="mdi-arrow-right-bold-box-outline"
              prev-icon="mdi-arrow-left-bold-box-outline"
              show-arrows
              color="deep-purple-accent-4"
              center-active
            >
              <b>
                <v-tab @click="oldtab = 0" key="-1" color="orange"> Favoris </v-tab>
              </b>
              <b v-for="cat in categories" :key="cat.categoryid">
                <v-tab
                  color="orange"
                  @click="
                    addcatbutton = false;
                    selectedcategory = cat;
                    oldtab = tab;
                  "
                  :key="cat.categoryid"
                  >{{ cat.categoryname }}

                  <v-icon
                    icon="mdi-menu-down"
                    style="padding: 6px"
                    size="x-large"
                    v-if="isadmin"
                  >
                  </v-icon>
                  <v-menu activator="parent">
                    <v-list v-if="isadmin">
                      <v-list-item>
                        <v-list-item-title @click="deleteCategory(cat)">
                          <v-icon
                            icon="mdi-trash-can"
                            style="margin-right: 15px"
                          ></v-icon>
                          Supprimer</v-list-item-title
                        >
                      </v-list-item>
                      <v-list-item>
                        <v-list-item-title @click="editCategory(cat)"
                          ><v-icon
                            icon="mdi-pencil"
                            style="margin-right: 15px"
                          ></v-icon
                          >Editer</v-list-item-title
                        >
                      </v-list-item>
                    </v-list>
                  </v-menu>
                </v-tab>
              </b>
              <b v-if="isadmin">
                <v-tab @click="openAddCategoryDialog()">
                  <v-icon size="32">mdi-plus</v-icon></v-tab
                >
              </b>
            </v-tabs>

            <v-tabs-window v-model="tab" id="menuwindow">
              <v-tabs-window-item v-for="n in 24" :key="n" :value="n">
                <v-container fluid style="text-align: center">
    <v-row ref="sortableContainer" class="sortable-grid">
      <v-col
        v-for="(element, index) in filtredMenu"
        :key="element.text"
        cols="12"
        md="2"
        class="sortable-item"
        :data-index="index"
      >
        <v-card
          :disabled="getMenuDisabled(btnvalidatedisabled, ticketvalidate, sessionsmap.get(loggeduser?.Username), isadmin, managerauthorisations, selectedwaiter)"
          style="background-color: transparent; line-height: 3mm"
          @click="elementClick(element)"
          :elevation="5"
          class="card-element"
          rounded="lg"
        >
        <v-menu activator="parent">
            <v-list v-if="isadmin">
              <v-list-item>
                <v-list-item-title
                  @click="deleteProduct(element)"
                  style="font-size: 14px"
                >
                  <v-icon
                    icon="mdi-trash-can"
                    style="margin-right: 5px; font-size: 14px; color: orangered"
                  ></v-icon>
                  Supprimer
                </v-list-item-title>
              </v-list-item>
              <v-list-item>
                <v-list-item-title
                  @click="editelement = element;openEditProductDialog();"
                  style="font-size: 14px"
                  
                >
                  <v-icon
                    icon="mdi-pencil"
                    style="margin-right: 5px; font-size: 14px; color: darkgreen"
                  ></v-icon>
                  Editer
                </v-list-item-title>
               
              </v-list-item>
              <v-list-item v-if="selectedcategory.categoryid>0" >
                <v-list-item-title
                 
                  style="font-size: 14px"
                >
                  <v-icon
                    icon="mdi-pencil"
                    style="margin-right: 5px; font-size: 14px; color: darkgreen"
                  ></v-icon>
                 Catégorie
                </v-list-item-title>
                <v-menu  activator="parent" submenu location="end">
                    <v-list>
                      <v-list-item v-for="cat in categories" :key="cat.categoryid" link @click="changeProductCategory(element, cat, selectedcategory.categoryid); ">
                        <v-list-item-title :class="element.category== cat.categoryid ? 'highlighted' : null">{{cat.categoryname}}</v-list-item-title>
                      </v-list-item>
                    </v-list>
                   
                  </v-menu>
              </v-list-item>
            </v-list>
          </v-menu>
          <v-img
            :src="getImageURL(element.image)"
            class="card-element-image"
            rounded="lg"
            cover
          ></v-img>
          <label class="card-label" style="font-size: 11px">{{ element.text }}
            
          </label>
          
        </v-card>
      </v-col>
      <v-col>
      <div class="d-flex flex-wrap ga-2">
        <v-card
          v-if="isadmin && tab != 0"
          @click="openAddProductDialog"
          style="background-color: transparent; line-height: 3mm; vertical-align: center"
          :elevation="10"
          class="card-element-no-drag"
          
          rounded="lg"
        >
          <v-icon size="48" style="margin-top: 50px">mdi-plus</v-icon>
        </v-card>
      </div>
    </v-col>
    </v-row>

    
    
  </v-container>
              </v-tabs-window-item>
            </v-tabs-window>
          </v-card>
          <v-spacer></v-spacer>
          <v-card style="float: right" class="mx-auto ticket" width="34%">
            <v-list-subheader style="font-weight: bolder; width: 100%">
              Ticket N° : {{ ticketnumber }}
            </v-list-subheader>
            <v-list-subheader style="font-weight: bolder; width: 100%" v-if="subscriberInfos?.name" >
              Nom de l'abonné(e) : {{ subscriberInfos.name }}
              <v-icon v-if="subscriberInfos?.isGold" color="amber-darken-2" size="small" class="ml-2">mdi-crown</v-icon>
            </v-list-subheader>
            <v-list style="max-height: 40vh; overflow-y: auto" id="listticket">
              <v-divider></v-divider>

              <v-list-item v-for="e in ticket[0].content" :key="e.ID">
                <v-speed-dial
                  :disabled="btnvalidatedisabled && !ticketvalidate"
                  location="left center"
                  transition="fade-transition"
                >
                  <template v-slot:activator="{ props: activatorProps }">
                    <v-chip
                      v-bind="activatorProps"
                      size="small"
                      color="orange"
                      :text="e.count"
                    ></v-chip>
                  </template>
                  <v-btn
                    key="1"
                    icon="mdi-comment-text"
                    style="color:mediumslateblue; background-color: yellow;"
                    size="small"
                    @click="showComment(e)"
                  ></v-btn>
                  <v-btn
                    key="1"
                    icon="$clear"
                    color="red"
                    size="small"
                    @click="removeFromTicket(e.id)"
                  ></v-btn>
                  <v-btn
                    key="2"
                    icon="$plus"
                    color="green"
                    size="small"
                    @click="increaseCount(e)"
                  ></v-btn>
                  <v-btn
                    key="3"
                    icon="$minus"
                    color="orange"
                    size="small"
                    style="color: white !important"
                    @click="reduceCount(e)"
                  ></v-btn>
                  <v-btn
                    key="4"
                    icon="mdi-dialpad"
                    color="orangered"
                    size="small"
                    style="color: white !important; background-color:grey !important;"
                    @click="showPad(e)"
                  ></v-btn>
                </v-speed-dial>
                <v-list-item-text style="font-size: 13px">
                  <label style="margin-left: 10px">{{ e.text }}</label>
                  <span
                    size="small"
                    style="
                      float: right;
                      margin-right: 10px;
                      align-content: center;
                    "
                    variant="elevated"
                    color="primary"
                    >{{ (e.count * parseFloat(e.price)).toFixed(2) }}</span
                  >
                </v-list-item-text>
              </v-list-item>
            </v-list>
            <v-divider></v-divider>
            <div style="height: 30px; margin-top: 5px">
              <label style="float: left; padding-left: 5px">Total : </label>
              <span
                style="
                  float: right;
                  margin-right: 10px;
                  color: whitesmoke !important  ;
                  font-size: 16px;
                "
                variant="elevated"
                color="orange"
                >{{ ticket[0].totalticket.toFixed(2) }}</span
              >
            </div>
            <v-divider></v-divider>
            
            <div
              style="
                height: 140px;
                margin-top: 5px;
                margin-left: 10px;
                margin-right: 10px;
                margin-bottom: 5px;
              "
            >
              <v-btn
                style="background-color: chocolate; float: left; width: 45%"
                v-on:click="emptyTicket(); emptyTicketTwo()"
                :disabled="btnvalidatedisabled && !ticketvalidate"
              >
                vider</v-btn
              >
              <v-btn
                style="background-color: cadetblue; float: right; width: 45%"
                v-on:click="validateTicket(sessionsmap, securitysettings, authOff)"
                :disabled="getBtnState(btnvalidatedisabled)"
              >
                Valider</v-btn
              >
              <v-btn
              v-on:click="printReceipt(securitysettings, authOff)"
              style="
                  background-color: darkolivegreen;
                  float: right;
                  width: 45%;
                  margin-top: 5px;
                "
                :disabled="ticketvalidate"
                
              >
                Imprimer</v-btn
              >
              <p></p>
              <v-btn
                style="
                  background-color: darkseagreen;
                  float: right;
                  width: 100%;
                  margin-top: 20px;
                  float: bottom;
                "
                :disabled="!btnvalidatedisabled"
                v-on:click="newTicket(); newTicketTwo()"
                id="newticketbutton"
                
              >
                Nouveau ticket</v-btn
              >
            </div>
            <div style="width: 99%; padding-left: 5px; padding-right: 5px">
              <v-select
                label="Serveur"
                :items="waitersnames"
                variant="outlined"
                ref="waiterbox"
                :disabled="getDisabledState()"
                v-model="selectedwaiter"
              ></v-select>
            </div>
          </v-card>
        </v-container>
      </v-main>
      <v-bottom-navigation >
                                                
        <v-btn value="Administrer" @click="showadminlogin(authOff)"  v-if="loggeduser?.Role=='Gérant' || loggeduser?.Username=='SETUP'" class="bottomnavbutton">
          <v-icon ref="adminicon" :class="adminclass">{{
            adminloginicon
          }}</v-icon>
          <span s>Admin</span>
        </v-btn>
        <v-btn value="session" @click="changeSessionStatus()" class="bottomnavbutton">
          <v-icon ref="session_icon" :class="sessionstatusclass">{{
            sessionstatusicon
          }}</v-icon>
         
          <span >{{sessionstatustext}}</span>
        </v-btn>
        <v-btn value="orderslist" @click="showOrdersList()" class="bottomnavbutton">
          <v-icon ref="orderslist" >mdi-format-list-bulleted</v-icon>
          <span >Commandes</span>
        </v-btn>

        <v-btn value="Situation" @click="showSituation()" class="bottomnavbutton">
          <v-icon  >mdi-cash</v-icon>
          <span >Situation</span>
        </v-btn>
        <v-btn value="vkeyboard" @click="showVirtualkeyBoard()" class="NOMOBILE" >
          <v-icon>mdi-keyboard</v-icon>

          <span >Clavier virtuel</span>
        </v-btn>
       

        <v-btn
          value="Ticket"
          @click="showTicketOnLowResolution()"
          class="showinlowresolution bottomnavbutton"
        >
          <v-icon>mdi-receipt</v-icon>

          <span>Ticket</span>
        </v-btn>
       
      </v-bottom-navigation>
      <v-card>
        <v-navigation-drawer
          v-if="isadmin"
          v-model="drawer"
          :location="$vuetify.display.mobile ? 'bottom' : undefined"
          temporary
        >
          <v-list density="compact" nav>
            <v-list-item
              v-if="isadmin"
              prepend-icon="mdi-tune"
              title="Paramètres"
              value="params"
              @click="
                paramsdialog = true;
                drawer = false;
              "
            ></v-list-item>

            <v-list-item
              v-if="isadmin"
              prepend-icon="mdi-account-group-outline"
              title="Gestion des utilisateurs"
              value="usersmanagement"
              @click="userslistsheet = !userslistsheet"
            ></v-list-item>
            <v-list-item
              v-if="isadmin"
              prepend-icon="mdi-printer-settings"
              title="Gestion de l'impression"
              value="printingmanagement"
              @click="printingmanagementsheet = ! printingmanagementsheet"
            ></v-list-item>
          </v-list>
        </v-navigation-drawer>
      </v-card>
      <div class="text-center pa-4">
        <v-dialog v-model="paramsdialog" width="auto" max-height="600" persistent>
          <v-card
            max-width="400"
            max-height="600"
            prepend-icon="mdi-tune"
            text="Ici vous pouvez modifier les paramètres de l'application"
            title="Pramètres de l'application"
          >
            <div>
              <v-expansion-panels>
                <v-expansion-panel>
                  <v-expansion-panel-title
                    collapse-icon="mdi-minus"
                    expand-icon="mdi-plus"
                  >
                    Pramètres généraux
                  </v-expansion-panel-title>
                  <v-expansion-panel-text>
                    
                    <v-text-field clearable v-model="generalsettings.posname" label="Intitulé du point de vent" variant="solo" dense  @focus="onUpdatefocused"></v-text-field>
                    <v-text-field clearable v-model="generalsettings.address" label="Adresse" variant="solo" dense  @focus="onUpdatefocused"></v-text-field>
                    <v-text-field clearable label="Téléphone" v-model="generalsettings.phone" variant="solo" dense  @focus="onUpdatefocused"></v-text-field>
                    <v-text-field clearable label="Code WIFI" v-model="generalsettings.wificode" variant="solo" dense  @focus="onUpdatefocused"></v-text-field>
                    <v-btn text="Enregister" v-on:click="saveGeneralInformations()" style="color:darkcyan; margin-left:65%" ></v-btn>
                    
                 </v-expansion-panel-text>
                </v-expansion-panel>
                <v-expansion-panel>
                  <v-expansion-panel-title
                    collapse-icon="mdi-minus"
                    expand-icon="mdi-plus"
                  >
                    Sécurité
                  </v-expansion-panel-title>
                  <v-expansion-panel-text>
                    
                   <v-radio-group label="Mode de déconnexion" v-model="securitysettings.disconnectmode">
                    <v-radio label="Déconnecter après validation" value="one"></v-radio>
                    <v-radio label="Déconnecter après impression" value="two"></v-radio>
                    
                  </v-radio-group>
                          
                        
                  <v-btn text="Enregister" v-on:click="saveSecurityOptions()" style="color:darkcyan; margin-left:65%"></v-btn>
                    
                 </v-expansion-panel-text>
                </v-expansion-panel>
                <v-expansion-panel>
                  <v-expansion-panel-title
                    collapse-icon="mdi-minus"
                    expand-icon="mdi-plus"
                  >
                    Autorisations
                  </v-expansion-panel-title>
                    <v-expansion-panel-text >
                          <v-checkbox
                          
                            v-model="managerauthorisations.managercancreatetickets"
                            label="Permettre au manager de créer des commandes."
                          ></v-checkbox>
                          <v-checkbox
                            
                            v-model="managerauthorisations.managercanselectwaiter"
                            label="Permettre au manager de selectionner le serveur manuellement."
                          ></v-checkbox>
                          <v-btn text="Enregister" v-on:click="saveAuthSettings()" style="color:darkcyan; margin-left:65%"></v-btn>
                    
                    </v-expansion-panel-text>
                </v-expansion-panel>
               
              </v-expansion-panels>
            </div>

            <template v-slot:actions>
              <v-btn
                class="ms-auto"
                text="Fermer"
                @click="paramsdialog = false"
              ></v-btn>
            </template>
          </v-card>
        </v-dialog>
      </div>
      <div class="text-center">
        <v-bottom-sheet v-model="sheet" inset>
          <v-card class="text-left" height="400">
            <v-btn size="x-large" @click="sheet = !sheet">Fermer</v-btn>

            <v-card-text>
           

              <div>
                <v-list-subheader style="font-weight: bolder; width: 100%">
                  Ticket N° : {{ ticketnumber }}
                </v-list-subheader>
                <v-list
                  style="max-height: 50vh; overflow-y: auto"
                  id="listticket"
                >
                  <v-divider></v-divider>

                  <v-list-item v-for="e in ticket[0].content" :key="e.ID">
                    <v-speed-dial
                      location="left center"
                      transition="fade-transition"
                    >
                      <template v-slot:activator="{ props: activatorProps }">
                        <v-chip
                          v-bind="activatorProps"
                          style=""
                          size="small"
                          color="orange"
                          :text="e.count"
                        ></v-chip>
                      </template>
                      <v-btn
                    key="1"
                    icon="mdi-comment-text"
                    style="color:mediumslateblue; background-color: yellow;"
                    size="small"
                    @click="showComment(e)"
                  ></v-btn>
                      <v-btn
                        key="1"
                        icon="$clear"
                        color="red"
                        size="small"
                        @click="removeFromTicket(e.id)"
                      ></v-btn>
                      <v-btn
                        key="2"
                        icon="$plus"
                        color="green"
                        size="small"
                        @click="increaseCount(e)"
                      ></v-btn>
                      <v-btn
                        key="3"
                        icon="$minus"
                        color="orange"
                        size="small"
                        style="color: white !important"
                        @click="reduceCount(e)"
                      ></v-btn>
                      <v-btn
                        key="4"
                        icon="mdi-dialpad"
                        color="orange"
                        size="small"
                        style="color: white !important; background-color : grey !important;"
                        @click="showPad(e)"
                      ></v-btn>
                    </v-speed-dial>
                    <v-list-item-text style="font-size: 11px">
                      <label style="margin-left: 5px">{{ e.text }}</label>
                      <span
                        size="small"
                        style="
                          float: right;
                          margin-right: 5px;
                          align-content: center;
                        "
                        variant="elevated"
                        color="primary"
                        >{{ (e.count * parseFloat(e.price)).toFixed(2) }}</span
                      >
                    </v-list-item-text>
                  </v-list-item>
                </v-list>
                <v-divider></v-divider>
                <div style="height: 30px; margin-top: 5px">
                  <label style="float: left; padding-left: 5px">Total : </label>
                  <span
                    size="small"
                    style="
                      float: right;
                      width: 120px;
                      padding: 0 25px;
                      text-align: right;
                      color: whitesmoke !important  ;
                      font-size: 14px;
                    "
                    variant="elevated"
                    color="orange"
                    >{{ ticket[0].totalticket.toFixed(2) }}</span
                  >
                </div>
                <v-divider></v-divider>
                <div
                  style="
                    height: 140px;
                    margin-top: 5px;
                    margin-left: 10px;
                    margin-right: 10px;
                    margin-bottom: 5px;
                  "
                >
                  <v-btn
                    style="background-color: orange; float: left; width: 45%"
                    v-on:click="emptyTicket(); emptyTicketTwo()"
                    :disabled="btnvalidatedisabled && !ticketvalidate"
                  >
                    vider</v-btn
                  >
                  <v-btn
                    style="
                      background-color: cadetblue;
                      float: right;
                      width: 45%;
                    "
                    v-on:click="validateTicket(sessionsmap,securitysettings, authOff)"
                    :disabled="btnvalidatedisabled"
                  >
                    Valider</v-btn
                  >
                  <v-btn
                   v-on:click="printReceipt(securitysettings, authOff)"
                    style="
                      background-color: darkolivegreen;
                      float: right;
                      width: 45%;
                      margin-top: 5px;
                    "
                    :disabled="ticketvalidate"
                  >
                    Imprimer</v-btn
                  >
                  <p></p>
                  <v-btn
                    style="
                      background-color: darkseagreen;
                      float: right;
                      width: 100%;
                      margin-top: 20px;
                      float: bottom;
                    "
                    :disabled="!btnvalidatedisabled"
                    v-on:click="newTicket(); newTicketTwo()"
                  >
                    Nouveau ticket</v-btn
                  >
                </div>
                <div style="width: 99%; padding-left: 5px; padding-right: 5px">
                  <v-select
                    label="Serveur"
                    :items="waitersnames"
                    variant="outlined"
                    ref="waiterbox"
                    v-model="selectedwaiter"
                    readonly
                  ></v-select>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-bottom-sheet>
      </div>

      <!--Add new Category dialog-->
      <new-category-dialog-component
        v-model="isCategoryDialogOpen"
        :categorytoedit="editcat"
        @alert="alertHandler"
        @update:model-value="closeCategoryDialog"
        :categories="categories"
      ></new-category-dialog-component>

      <!--Add new product dialog-->
      <new-product-dialog-component
        v-if="isadmin"
        :parentcategory="selectedcategory"
        v-model="isProductDialogOpen"
        @alert="alertHandler"
        @update:model-value="closeProductDialog"
      ></new-product-dialog-component>
      <!--Edit product dialog-->
      <edit-product-dialog-component
        v-if="isadmin"
        :editproduct="editelement"
        v-model="isEditProductDialogOpen"
        @alert="alertHandler"
        @update:model-value="closeEditProductDialog"
        :categories="categories"
      ></edit-product-dialog-component>
      <!--User management dialog-->
      <manage-users-dialog-component
        v-if="isadmin"
        v-model="isManageUsersDialogOpen"
        @alert="alertHandler"
        @update:model-value="closeManageUsersDialog"
        @inform="inform"
      ></manage-users-dialog-component>
      <!--Edit users dialog-->
      <edit-users-dialog-component
        v-if="isadmin"
        v-model="isEditUsersDialogOpen"
        @alert="alertHandler"
        @update:model-value="closeEditUsersDialog"
        @inform="inform"
        :usertoedit="usertoedit"
      ></edit-users-dialog-component>
     
      
   

      <!--alert dialog-->
      <div class="centerstuff">
        <v-bottom-sheet v-model="showalerttwo" inset >
          <v-card
            :class="alertcardclasstwo"
            id="alertcard"
            height="150"
            style="opacity: 0.9"
          >
            <v-card-text>
              <v-btn variant="text" @click="showalerttwo = !showalerttwo">
                Fermer
              </v-btn>
              <div>
                {{ alerttexttwo }}
              </div>
            </v-card-text>
          </v-card>
        </v-bottom-sheet>
      </div>
      <div class="centerstuff">
        <v-bottom-sheet v-model="showalert" inset >
          <v-card
            :class="alertcardclass"
            id="alertcard"
            height="150"
            style="opacity: 0.9"
          >
            <v-card-text>
              <v-btn variant="text" @click="showalert = !showalert">
                Fermer
              </v-btn>
              <div>
                {{ alerttext }}
              </div>
            </v-card-text>
          </v-card>
        </v-bottom-sheet>
      </div>

      <div class="text-center pa-4">
        <v-dialog v-model="confirmdialog" max-width="400" persistent>
          <v-card
            prepend-icon="mdi-delete"
            :text="confirmationtext"
            title="Supression"
          >
            <template v-slot:actions>
              <v-spacer></v-spacer>

              <v-btn
                @click="
                  confirmDelete();
                  confirmdialog = false;
                "
              >
                OUI
              </v-btn>

              <v-btn @click="confirmdialog = false"> NON </v-btn>
            </template>
          </v-card>
        </v-dialog>
      </div>
      <!--COMMENT DIALOG-->
      <v-dialog v-model="showcommentdialog" max-width="400" persistent>

          
          <v-card
            prepend-icon="mdi-comment-text"
            
            title="Message"
          >
            <v-textarea
              class="ma-2"
              v-model="comment"
              label="Message"
              bg-color="amber-lighten-4"
           
              @focus="onUpdatefocused"
              clearable></v-textarea>
             
              <v-btn @click="saveComment(selecteditem);">Enregistrer</v-btn>
              <v-btn @click="showcommentdialog = false"> Fermer </v-btn>
          </v-card>
        </v-dialog>
        <!--DELETE USER CONFIRMATION-->
      <div class="text-center pa-4">
        <v-dialog v-model="confirmdeleteprintingprofiledialog" max-width="400" persistent>
          <v-card
            prepend-icon="mdi-delete"
            :text="confirmtext"
            title="Supression de profile d'impression"
          >
            <template v-slot:actions>
              <v-spacer></v-spacer>

              <v-btn
                @click="
                  deletePrintingProfile();
                  confirmdeleteprintingprofiledialog = false;
                "
              >
                OUI
              </v-btn>

              <v-btn @click="confirmdeleteuderdialog = false"> NON </v-btn>
            </template>
          </v-card>
        </v-dialog>
      </div>
      <!--DELETE USER CONFIRMATION-->
      <div class="text-center pa-4">
        <v-dialog v-model="confirmdeleteuderdialog" max-width="400" persistent>
          <v-card
            prepend-icon="mdi-delete"
            :text="confirmtext"
            title="Supression"
          >
            <template v-slot:actions>
              <v-spacer></v-spacer>

              <v-btn
                @click="
                  deleteUser();
                  confirmdialog = false;
                "
              >
                OUI
              </v-btn>

              <v-btn @click="confirmdeleteuderdialog = false"> NON </v-btn>
            </template>
          </v-card>
        </v-dialog>
      </div>
       <!--CONFIRM CHECK ORDER-->
       <div class="text-center pa-4">
        <v-dialog
          v-model="confirmcheckorderdialog"
          max-width="400"
          persistent
        >
          <v-card
            prepend-icon="mdi-marker-check"
            :text="confirmcheckordertext"
            title="Paiement"
          >
            <template v-slot:actions>
              <v-spacer></v-spacer>

              <v-btn
                @click="
                  checkConfirmation();
                  confirmcheckorderdialog= false;
                "
              >
                OUI
              </v-btn>

              <v-btn @click="confirmcheckorderdialog= false"> NON </v-btn>
            </template>
          </v-card>
        </v-dialog>
      </div>
      <!--CONFIRM CANCEL ORDER-->
      <div class="text-center pa-4">
        <v-dialog
          v-model="confirmcancelorderdialg"
          max-width="400"
          persistent
        >
          <v-card
            prepend-icon="mdi-cancel"
            :text="confirmcancelordertext"
            title="Annulation"
          >
            <template v-slot:actions>
              <v-spacer></v-spacer>

              <v-btn
                @click="
                  cancelOrder();
                  confirmcancelorderdialg= false;
                "
              >
                OUI
              </v-btn>

              <v-btn @click="confirmcancelorderdialg= false"> NON </v-btn>
            </template>
          </v-card>
        </v-dialog>
      </div>
      <!--CONFIRM START END SESSION-->
      <div class="text-center pa-4">
        <v-dialog
          v-model="confirmstartendsessiondialog"
          max-width="400"
          persistent
        >
          <v-card
            prepend-icon="mdi-clock-start"
            :text="confirmstartendsessiontext"
            title="Session"
          >
            <template v-slot:actions>
              <v-spacer></v-spacer>

              <v-btn
                @click="
                  confirmStartOrEndSession();
                  confirmstartendsessiondialog = false;
                "
              >
                OUI
              </v-btn>

              <v-btn @click="confirmstartendsessiondialog = false"> NON </v-btn>
            </template>
          </v-card>
        </v-dialog>
      </div>
      <div class="text-center pa-4">
        <v-dialog
          v-model="deleteproductconfirmdialog"
          max-width="400"
          persistent
        >
          <v-card
            prepend-icon="mdi-delete"
            :text="confirmationtext"
            title="Supression"
          >
            <template v-slot:actions>
              <v-spacer></v-spacer>

              <v-btn
                @click="
                  confirmDeleteProduct();
                  deleteproductconfirmdialog = false;
                "
              >
                OUI
              </v-btn>

              <v-btn @click="deleteproductconfirmdialog = false"> NON </v-btn>
            </template>
          </v-card>
        </v-dialog>
      </div>
      <!--KEYPAD-->
      <div>
        <v-dialog v-model="showpad" width="150px" class="centerstuff">
          <div class="display" id="display">{{ currentelement.count }}</div>
          <div class="keypad">
            <button @click="pressButton('1')">1</button>
            <button @click="pressButton('2')">2</button>
            <button @click="pressButton('3')">3</button>
            <button @click="pressButton('4')">4</button>
            <button @click="pressButton('5')">5</button>
            <button @click="pressButton('6')">6</button>
            <button @click="pressButton('7')">7</button>
            <button @click="pressButton('8')">8</button>
            <button @click="pressButton('9')">9</button>
            <button class="control" @click="clearDisplay()">C</button>
            <button @click="pressButton('0')">0</button>
            <button class="control" @click="deleteLast()">⌫</button>
            <button class="validate" @click="validateInput()">✔</button>
          </div>
        </v-dialog>
      </div>
      <!--PRINTING DIALOG-->

      <v-dialog v-model="showprintingdialog" width="400" height="700" class="centerstuff">
          
        
            <v-card
              prepend-icon="mdi-printer"
              class="pa-2"
             
              title="Impression"
            >
                    <v-text-field
                label="Nom du profile"
                v-model="printingprofile.profile"
                variant="solo-filled"
                :rules="[required]"
                density="compact"
                
              ></v-text-field>
              <v-select
                label="Mode d'impression"
                v-model="printingprofile.printingmode"
                :items="printingmodes"
                variant="solo-filled"
                density="compact"
                :rules="[required]"
              >
              </v-select>
              <v-select
                label="Fonction"
                v-model="printingprofile.printingfunction"
                :items="printingfunctions"
                variant="solo-filled"
                density="compact"
                :rules="[required]"
              >
               
              </v-select>
              <v-card v-if ="printingprofile.printingfunction=='Personnalisé'" style="overflow-y: scroll;" min-height="80" >
                <v-chip-group
                        v-model="printingprofile.printingoptions"
                        column
                        multiple
                      >
                        <v-chip v-for="cat in categories"
                        style="font-size: 0.6rem; "
                          :key="cat.categoryid"
                          :text="cat.categoryname"
                          :value="cat.categoryid"
                          variant="outlined"
                          filter
                        ></v-chip>
                </v-chip-group>
              </v-card>
              <v-text-field
                label="Adresse IP"
                v-model="printingprofile.printerip"
                variant="solo-filled"
                density="compact"
                :rules="[required]"
                      >
              </v-text-field>
              <v-text-field
                label="Port"
                v-model="printingprofile.printerport"
                variant="solo-filled"
                density="compact"
                :rules="[required]"
                >
              </v-text-field>
              <v-number-input 
                v-model="printingprofile.copies"
                label="Copies"  
                density="compact"
                variant="solo-filled"
                control-variant="hidden"
                :rules="[required]"
                min="1">
                
                
             </v-number-input>
             <v-switch
              v-model="printingprofile.active"
              color="green"
              :label="printingprofile.active ? 'Actif' : 'Inactif'"
              
              hide-details
            ></v-switch>

              

                
               
            
                    <template v-slot:actions>
                <v-spacer></v-spacer>
  
               
  
                <v-btn @click="showprintingdialog = false"> Fermer </v-btn>
                <v-btn @click="savePrintingProfile" :disabled="!Object.values(printingprofile).every(value => value !== '' && value !== null && value !== undefined)" > Enregistrer </v-btn>
              </template>
            </v-card>
         
        </v-dialog>
      
       
 
      <!--ADMIN LOGIN-->
     
      <div>
        <v-dialog v-model="showadminprompt" width="150px" class="centerstuff">
          <div style="border-radius: 10px; width: 150px; text-align: center">
            <input
              ref="adminpass"
              type="password"
              v-model="adminpassword"
              style="
                width: 80%;
                color: black;
                text-align: center;
                font-size: 20px;
                border-color: darkslateblue;
                background-color: darkseagreen;
                border-radius: 10px;
                border-style: solid;
                border-width: 3px;
              "
               @focus="onUpdatefocused"
            />
            <v-btn
              style="
                margin-top: 20px;
                background-color: darkolivegreen;
                color: white;
              "
              @click="adminLogin()"
              >Go</v-btn
            >
          </div>
        </v-dialog>
      </div>
      <!--ORDERS LIST-->
      <div>
      
      <v-bottom-sheet v-model="showorderslist" >
        <div><order-content-component
       :order="selectedorder"
       :settings="generalsettings"
       v-model="showcontent"
       @close="showcontent=false"
     ></order-content-component></div>
        
        <v-card class="text-center" height="100vh">
          <v-card-text>
            <v-btn variant="text" @click="showorderslist=!showorderslist">
              Fermer
            </v-btn>

            <div>
              <v-card flat>
           
                <v-card-title class="d-flex">
                 
                    
                      <input type="date" v-model="startdatemodel" >
                      
                      <label style="font-size: 12px; padding-top: 15px; margin-right: 10px; margin-left:10px; width:100px; align-content: right;">à partir de :</label>
                      <input type="time" id="appt-time" name="appt-time"  step="600" value="00:00:01" v-model="starttimemodel">
                      <label for="status" style="font-size: 12px; padding-top: 15px; margin-right: 10px; margin-left:10px">Etat:</label>
                      <select id="statusbox" name="language" v-model="ordersstatus" @change="statuschanged">
                          <option value="TOUT" selected>Tout</option>
                          <option value="IMPAYE">Impayé</option>
                          <option value="PAYE">Payé</option>
                          <option value="ANNULE">Annulé</option>
                      </select>
                      
                      <v-spacer></v-spacer>
                     
                      <v-text-field
                        v-model="searchorder"
                        density="compact"
                        label="Recherche"
                        prepend-inner-icon="mdi-magnify"
                        variant="solo-filled"
                        flat
                        hide-details
                        single-line
                        max-width="300"
                        style="height: 40px"
                         @focus="onUpdatefocused"
                      ></v-text-field>
                  
                      
                 
                 
                </v-card-title>
                <v-card-title class="d-flex">
                 
                    
                 <input type="date" v-model="enddatemodel" >
                 
                 <label style="font-size: 12px; padding-top: 15px; margin-right: 10px; margin-left:10px; width:100px; align-content: right;">jusqu'à :</label>
                 <input type="time" id="appt-time" name="appt-time" step="600" value="23:59:59" v-model="endtimemodel">
                 
                 <v-btn style="margin-top: 5px;margin-left: 10px;" @click="getOrders"><v-icon>mdi-refresh</v-icon> Recharger</v-btn>
                 <v-btn style="margin-top: 5px;margin-left: 10px;"  @click="printOrders()"><v-icon>mdi-file-pdf</v-icon> Exporter</v-btn>
                 <v-spacer></v-spacer>
                
             
                 
            
            
           </v-card-title>




                <v-data-table
                  id="orderstable"
                  :headers="[
                    {
                      title: 'Numéro',
                      align: 'start',
                      sortable: true,
                      key: 'number',
                    },
                    { title: 'Serveur', key: 'waiter' },

                    {title:'Initiateur', key:'creator'},
                    { title: 'Date et heure', key: 'creationdate' },
                    {title:'Etat', key:'status'},
                    {title:'Total', key:'totalticket'},
                    { title: 'Actions', key: 'actions' },
                  ]"
                  v-model:search="searchorder"
                  :filter-keys="['waiter', 'creator', 'status','number']"
                  :items="filtredorders"
                  style="text-align: left"
                   :itemsPerPageText="ordersProps.itemsPerPageText"
                   :itemsPerPageOptions="ordersProps.itemsPerPageOptions"
                  
                >
      
                <template v-slot:[`item.status`]="{item}">
                  <p :class="item.status">{{item.status}}</p>
                </template>
                <template v-slot:[`item.totalticket`]="{item}">
                  <p >{{item.totalticket.toFixed(2)}}</p>
                </template>
                  <template v-slot:[`item.actions`]="{ item }">
                    <v-icon
                      class="me-2"
                      size="large"
                      @click="confirmCheckOrder(item)"
                    >
                    mdi-cash-check

                    </v-icon>

                    <v-icon size="large" style="margin-left:5px" @click="showContent(item)">
                      mdi-eye
                    </v-icon>
                    <v-icon size="large" style="margin-left:5px" @click="confirmCancelOrder(item)" v-if="loggeduser?.Role=='Gérant'">
                      mdi-cancel
                    </v-icon>
                  </template>
              
                </v-data-table>
              </v-card>
            </div>
          </v-card-text>
        </v-card>
        
      </v-bottom-sheet>
    </div>
    <!--SITUATION-->
    <div>
      
      <v-bottom-sheet v-model="showsituation" inset  class="situation">
        <v-card class="text-center" height="55vh">
          <v-card-text>
            <v-btn variant="text" @click="showsituation = !showsituation">
              Fermer
            </v-btn>
           
         
              <v-card flat>
                <v-card-title class="d-flex align-center pe-2">
                 
                  Situation des commandes

                  <v-spacer></v-spacer>

                 
                </v-card-title>

                <v-divider></v-divider>
                <v-spacer></v-spacer>
                <v-select
                    label="Utilisateur"
                    :items="usersnames"
                    variant="outlined"
                    @update:modelValue="userSelectionChanged"
                    v-model="selectedperson"
                    
                  ></v-select>
                  <div style="font-size: 12px; text-align: left;">
                    <table style="width: 100%; border-collapse: collapse;">
                      <tr>
                        <td style="width: 60%;">Début  </td>
                        
                        <td style="width: 40%; text-align: right;">{{ getDateTime(sessions.get(selectedperson)?.start) }}</td>
                      </tr>
                      <tr >
                        <td style="width: 60%;">Fin  </td>
                        
                        <td style="width: 40%; text-align: right;">{{ getDateTime(sessions.get(selectedperson)?.end) }}</td>
                      </tr>
                      <tr >
                        <td style="width: 60%; background-color: darkcyan;">Total payées  </td>
                        
                        <td style="width: 40%; text-align: right; background-color: darkcyan;">{{getDateTime(sessions.get(selectedperson)?.start)!="" ? situation?.totalpaid?.toFixed(2): "0.00"}}</td>
                      </tr>
                      <tr>
                        <td style="width: 60%; background-color: darkcyan">Nombre payées  </td>
                        
                        <td style="width: 40%; text-align: right; background-color: darkcyan">{{getDateTime(sessions.get(selectedperson)?.start)!="" ? situation?.countpaid:"0.00"}}</td>
                      </tr>
                      <tr>
                        <td style="width: 60%; background-color: darkgoldenrod">Total impayées  </td>
                        
                        <td style="width: 40%; text-align: right;background-color: darkgoldenrod">{{getDateTime(sessions.get(selectedperson)?.start)!="" ? situation?.totalinpaid?.toFixed(2):"0.00"}}</td>
                      </tr>
                      <tr>
                        <td style="width: 60%;background-color: darkgoldenrod;">Nombre impayées  </td>
                        
                        <td style="width: 40%; text-align: right;background-color: darkgoldenrod">{{getDateTime(sessions.get(selectedperson)?.start)!="" ? situation?.countinpaid: "0.00"}}</td>
                      </tr>
                      <tr>
                        <td style="width: 60%;background-color: dimgray;">Total annulées  </td>
                        
                        <td style="width: 40%; text-align: right;background-color: dimgray">{{getDateTime(sessions.get(selectedperson)?.start)!="" ? situation?.totalcanceled?.toFixed(2):"0.00"}}</td>
                      </tr>
                      <tr>
                        <td style="width: 60%;background-color: dimgray">Nombre annulées  </td>
                        
                        <td style="width: 40%; text-align: right;background-color: dimgray">{{getDateTime(sessions.get(selectedperson)?.start)!="" ? situation?.countcanceled:"0.00"}}</td>
                      </tr>
                      <tr>
                        <td style="width: 60%;background-color:darkgreen; font-size: 14px; font-weight: bold;">Total  </td>
                        
                        <td style="width: 40%; text-align: right;background-color: darkgreen;font-size: 14px; font-weight: bold ">{{getDateTime(sessions.get(selectedperson)?.start)!="" ? situation?.total?.toFixed(2):"0.00"}}</td>
                      </tr>
                      
                      
                    </table>
                    <v-btn color="orange" style="float:right; margin-top:10px" @click="printSituation()">Imprimer</v-btn>
                  </div>

                
              </v-card>
            
          </v-card-text>
        </v-card>
      </v-bottom-sheet>
    </div>
    <!--PRINTING MANAGEMENT-->
      <div>
        <v-bottom-sheet v-model="printingmanagementsheet" inset>
          <v-card class="text-center" height="70vh" width="80vw" >
            <v-card-text>
              <v-btn variant="text" @click="printingmanagementsheet= !printingmanagementsheet">
                Fermer
              </v-btn>
              <v-btn
                 
                  ariant="text"
                  
                  @click="openManagePrintingDialog()"
                >
                  Nouveau
                </v-btn>

            

              <div>
                <v-card flat>
                  <v-card-title class="d-flex align-center pe-2">
                    <v-icon icon="mdi-printer"></v-icon> &nbsp;
                    Gestion de l'impression

                    
                    
                   
                  </v-card-title>

                  <v-spacer></v-spacer>


                  <v-data-table
                    :headers="[
                      {
                        title: 'Profile',
                        align: 'start',
                        sortable: true,
                        key: 'profile',
                      },
                      { title: 'Mode', key: 'printingmode' },
                      { title: 'Fonction', key: 'printingfunction' },
                      { title: 'Adresse IP', key: 'printerip' },
                      { title: 'Port', key: 'printerport' },
                      { title: 'Copies', key: 'copies' },
                      { title: 'Actif', key: 'active' },
                      { title: 'Actions', key: 'actions' },
                    ]"
                    
                    :items="printingprofiles"
                    style="text-align: left"
                     :itemsPerPageText="footerPropsPrinters.itemsPerPageText"
                     :itemsPerPageOptions="footerPropsPrinters.itemsPerPageOptions"
                    
                  >
                  <template v-slot:[`item.active`]="{ item }">
                    <p :class="item.active ? 'active' : 'inactive'">{{item.active ? 'Oui' : 'Non'}}</p>
                  </template>
                    <template v-slot:[`item.actions`]="{ item }">
                      <v-icon
                        class="me-2"
                        size="small"
                        @click="
                          printingprofile = item;
                          openEditPrintingProfile();
                        "
                      >
                        mdi-pencil
                      </v-icon>

                      <v-icon size="small" @click="confirmDeletePrintingProfile(item)">
                        mdi-delete
                      </v-icon>
                    </template>
                   
                  </v-data-table>
                </v-card>
              </div>
            </v-card-text>
          </v-card>
        </v-bottom-sheet>
      </div>
      <!--USERS LIST-->
      <div>
      
        <v-bottom-sheet v-model="userslistsheet" inset>
          <v-card class="text-center" height="70vh">
            <v-card-text>
              <v-btn variant="text" @click="userslistsheet = !userslistsheet">
                Fermer
              </v-btn>
              <v-btn
                 
                  ariant="text"
                  
                  @click="openManageUsersDialog()"
                >
                  Nouveau
                </v-btn>

            

              <div>
                <v-card flat>
                  <v-card-title class="d-flex align-center pe-2">
                    <v-icon icon="mdi-account-group-outline"></v-icon> &nbsp;
                    Liste des utilisateurs

                    <v-spacer></v-spacer>
                    
                    <v-text-field
                      v-model="search"
                      density="compact"
                      label="Recherche"
                      prepend-inner-icon="mdi-magnify"
                      variant="solo-filled"
                      flat
                      hide-details
                      single-line
                       @focus="onUpdatefocused"
                    ></v-text-field>
                  </v-card-title>

                  <v-divider></v-divider>
                  <v-spacer></v-spacer>


                  <v-data-table
                    :headers="[
                      {
                        title: 'Nom d\'utilisateur',
                        align: 'start',
                        sortable: true,
                        key: 'Username',
                      },
                      { title: 'Rôle', key: 'Role' },
                      { title: 'Actions', key: 'actions' },
                    ]"
                    v-model:search="search"
                    :filter-keys="['Username', 'Role']"
                    :items="users"
                    style="text-align: left"
                     :itemsPerPageText="footerProps.itemsPerPageText"
                     :itemsPerPageOptions="footerProps.itemsPerPageOptions"
                    
                  >
                    <template v-slot:[`item.actions`]="{ item }">
                      <v-icon
                        class="me-2"
                        size="small"
                        @click="
                          usertoedit = item;
                          openEditUsersDialog();
                        "
                      >
                        mdi-pencil
                      </v-icon>

                      <v-icon size="small" @click="confirmDeleteUser(item)">
                        mdi-delete
                      </v-icon>
                    </template>
                  
                  </v-data-table>
                </v-card>
              </div>
            </v-card-text>
          </v-card>
        </v-bottom-sheet>
      </div>
 
    </v-app>
  </v-responsive>
</template>

<script setup>
import { filtredMenu, getSortedFiltredMenu, initSortable, store } from "./Utils";  
import { ref, onMounted, reactive, watch} from "vue";
import jsPDF from "jspdf";
import "jspdf-autotable";
const initialprintingprofile = {profile:'', printingmode:'', printingfunction:'', printingoptions:[], printerip:'', printerport:'9100', copies:1, active:false}
const printingprofile=ref({...initialprintingprofile})
const printingmodes = ref(["USB", "LAN"])
const printingfunctions = ref(["Reçu", "Situation","Personnalisé"])
const printingprofiles = ref([])
const showprintingdialog = ref(false)

const selectedperson = ref("")
const isCategoryDialogOpen = ref(false);
const isProductDialogOpen = ref(false);
const isEditProductDialogOpen = ref(false);
const isManageUsersDialogOpen = ref(false);
const isEditUsersDialogOpen = ref(false);
const confirmdeleteuderdialog = ref(false);
const confirmtext = ref("");
const managerauthorisations = ref({managercancreatetickets:true, managercanselectwaiter:false})
const selectedwaiter = ref("");
const rail = ref(true);
const drawer = ref(false);
const waiters = ref("");
const waitersnames = ref("");
const selecteditem = ref(null);
const comment = ref("");
const authenticated=ref(false)
const users = ref("");
const userslistsheet = ref(false);
const showorderslist= ref(false)
const search = ref("");
const searchorder = ref('')
const usertodelete = ref("");
const selectedorder=ref(null)
const showsituation = ref(false)
let usertoedit = null;
const loggeduser = ref(null)
const ordersstatus=ref("Tout")
const sessionended = ref(false)
const printingmanagementsheet =ref(false)

const footerProps = {
  itemsPerPageText: "Utilisateurs par page", 
  itemsPerPageOptions: [5, 10, 20] 
} 
const footerPropsPrinters = {
  itemsPerPageText: "Imprimantes par page", 
  itemsPerPageOptions: [5, 10, 20] 
}
const ordersProps = {
  itemsPerPageText: "Commmandes par page", 
  itemsPerPageOptions: [10,20,50,100] 
} 
const subscriberInfos = ref(null)
const orders = ref([])
const filtredorders =ref([])

let editcat = { categoryname: "", categoryid: "" };
let addcatbutton = false;
let keyssaved = "";
let lenghtcounter = 0;
const showcontent=ref(false)
const startdatemodel=ref(null)
const starttimemodel=ref("00:00:01")
const enddatemodel=ref(null)
const usersnames = ref('')
const endtimemodel=ref("23:59:59")
const theme= ref("dark")
const generalsettings = ref({posname:'', address:'', phone:'', wificode:''})
const securitysettings= ref({disconnectmode:'two'})
const waiterbox = ref(null);
const themes = ref([{title:"bluegrey"},{title:"dark"},{title:"brown"},{title:"teal"},{title:"deepblue"},  {title:"purple"}])
function openManagePrintingDialog(){
  showprintingdialog.value = true
  printingprofile.value = {...initialprintingprofile}
}
const editprofile = ref(false)
function openEditPrintingProfile(){
  editprofile.value = true
  showprintingdialog.value = true
}
const printingprofiletodelete = ref("")
function getOverlayStyle(){
  return "z-index:9998;width:100%; height:100%; position: fixed; background-size: cover;top: 0;left: 0; background-image:url(https://"+ getHostName() +"/img/coffeeshop.jpg); opacity:0.5"
}
function confirmDeletePrintingProfile(item){
  confirmtext.value = `Voulez-vous vraiment supprimer le profile d'impression ${item.profile} ?`
  confirmdeleteprintingprofiledialog.value = true
  printingprofiletodelete.value = item
}
const confirmdeleteprintingprofiledialog = ref(false)
const valdiationbuttonstate = ref(false)
const getBtnState = (state) => {
  valdiationbuttonstate.value = state;
  return state;
};
watch(valdiationbuttonstate, (newValue) => {
  if (newValue === false) {
   
    enableRFIDReading();
    
  } else {
    
    disableRFIDReading();
    
  }
});
async function savePrintingProfile(){
 
  if (printingprofiles.value.find(e=>e.profile==printingprofile.value.profile)!=null && !editprofile.value){
    alertHandlerTwo({Prefix:"ERROR",Message:"Ce nom de profile existe déjà", data:""})
    return
  }
  
  if (editprofile.value){
    await updatePrintingProfile()
    return
  }
  printingprofiles.value.push({...printingprofile.value})
  const form = new FormData();
  
  form.append("printingsettings", JSON.stringify(printingprofiles.value))
  await fetch("https://" +  getHostName() + "/saveprintingsettings",{method:'post', body:form})
  .then(resp=>{
    if (!resp.ok) throw Error(resp.statusText)
    return resp.json()
  })
  .then(data=>{
    if (data==null) return
    alertHandlerTwo(data)
    
    printingprofile.value = {...initialprintingprofile}
    showprintingdialog.value = false
  })
  .catch(ex=>{
    console.error(ex)
  })
}
async function deletePrintingProfile(){
  
  const form = new FormData();
  form.append("printingsettings", JSON.stringify(printingprofiles.value))
  await fetch("https://" +  getHostName() + "/saveprintingsettings",{method:'post', body:form})
  .then(resp=>{
    if (!resp.ok) throw Error(resp.statusText)
    return resp.json()
  })
  .then(data=>{
    if (data==null) return
    alertHandlerTwo(data)
    if (data.Prefix=="ERROR") return
    printingprofiles.value = printingprofiles.value.filter(e=>e.profile!=printingprofiletodelete.value.profile)
    confirmdeleteprintingprofiledialog.value = false
   
    
  })
  .catch(ex=>{
    console.error(ex)
  })
}
async function updatePrintingProfile(){
 
  const form = new FormData();
  form.append("printingsettings", JSON.stringify(printingprofiles.value))
  await fetch("https://" +  getHostName() + "/saveprintingsettings",{method:'post', body:form})
  .then(resp=>{
    if (!resp.ok) throw Error(resp.statusText)
    return resp.json()
  })
  .then(data=>{
    if (data==null) return
    alertHandlerTwo(data)
    
    showprintingdialog.value = false
    
  })
  .catch(ex=>{
    console.error(ex)
  })
}
function changeTheme(t){
   theme.value = t
   localStorage.setItem('theme', theme.value)
}
function convertDateFormat(dateString) {
  
  const [year, month, day] = dateString.split("-");

 
  return `${day}-${month}-${year}`;
}
async function changeProductCategory(element,category, catid){
  
  let form = new FormData()
  form.append("elementid", element.productid)
  form.append("categoryid", category.categoryid)
 await fetch("https://" + getHostName() + "/changeproductcategory", {method:'post', body: form})
 .then(r=>{
  if (!r.ok) throw Error(r.statusText)
  return r.json()
 })
 .then(data=>{
   if (data==null) return
   alertHandlerTwo(data)
   let localoriginfiltredmenu = JSON.parse(localStorage.getItem("filtredMenu" + catid))
   let localdestinationfiltredmenu = JSON.parse(localStorage.getItem("filtredMenu" + category.categoryid))
   
   
   element.category=category.categoryid
   filtredMenu.value = filtredMenu.value.filter(e=>e.productid!=element.productid)  
   if (localdestinationfiltredmenu!=null){
      localdestinationfiltredmenu.push(element)
   }
   if (localoriginfiltredmenu!=null){
    localoriginfiltredmenu = localoriginfiltredmenu.filter(e=>e.productid!=element.productid)
   }
   
   
  
   
   localStorage.setItem("filtredMenu" + catid , JSON.stringify(localoriginfiltredmenu));
   localStorage.setItem("filtredMenu" + category.categoryid , JSON.stringify(localdestinationfiltredmenu));
   
 })
 .catch(ex=>{
   console.error(ex)
 })
}
const formvalidate = ref(false)
const required = (v)=> {
        
        return !!v || 'Champ requis';
        
}
function printOrders(){
  const pdf = new jsPDF();
  let sum=0
  filtredorders.value.forEach(item=>{
    if (item.status != "ANNULE"){
      sum += item.totalticket
    }
  })
  
 

  const tableheaders=[
                    {
                      title: 'Numéro',
                      align: 'start',
                      sortable: true,
                      key: 'number',
                    },
                    { title: 'Serveur', key: 'waiter' },
                    {title:'Initiateur', key:'creator'},
                    { title: 'Date et heure', key: 'creationdate' },
                    {title:'Etat', key:'status'},
                    {title:'Total', key:'totalticket'},
                    
                  ]
const headers = tableheaders.map((header) => header.title);
const rows = filtredorders.value.map((item) =>
  tableheaders.map((header) => item[header.key])
);
const now = new Date();
const hours = now.getHours(); 
const minutes = now.getMinutes(); 
const currentTime = `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`;


pdf.setFontSize(18);
pdf.text("Etat des opérations", 10, 20)
pdf.setFontSize(12)
pdf.text("Date de début: " + convertDateFormat(startdatemodel.value) +" à " + starttimemodel.value, 10, 30)
pdf.text("Date de fin: " + convertDateFormat(enddatemodel.value) + " à " + currentTime, 10, 40)
pdf.autoTable({
        head: [headers], 
        body: rows,
        startY: 50, 
        margin: { left: 10, right: 10 }, 
        styles: { fontSize: 10 }, 
        footStyles: { fontStyle: "bold" }, 
});
const pageWidth = pdf.internal.pageSize.getWidth();
const tableEndY = pdf.autoTable.previous.finalY + 10;

pdf.setFontSize(14);
pdf.setFont("helvetica", "bold");
pdf.text(
  `Total: ${sum.toFixed(2)}`,
  pageWidth - 50,
  tableEndY
);
pdf.save("Etatopérations.pdf");
}
function getDisabledState(){
  if (managerauthorisations.value.managercanselectwaiter){
    if (loggeduser.value?.Role == "Gérant"){
      return false
    }
    return true
  }
  return true
}
async function savePrintingSettings(){
  const form = new FormData()
  
  form.append("printingsettings", JSON.stringify(printingSettings.value))
  await fetch("https://" + getHostName() + "/saveprintingsettings", {method:'post', body: form})
  .then(resp=>{
    if (!resp.ok) throw Error(resp.statusText)
    return resp.json()
  })
  .then(data=>{
    if (data==null) return
    
    alertHandlerTwo(data)
   
  })

  .catch(ex=>{
    console.error(ex)
  })
}
async function saveAuthSettings(){
  const form = new FormData()
  form.append("managersettings", JSON.stringify(managerauthorisations.value))
  await fetch("https://" + getHostName() + "/savemanagersettings", {method:'post', body: form})
  .then(resp=>{
    if (!resp.ok) throw Error(rersp.statusText)
    return resp.json()
  })
  .then(data=>{
    if (data==null) return
    
    alertHandlerTwo(data)
   
  })
  .catch(ex=>{
    console.error(ex)
  })
}
async function loadSettings(){

  await fetch("https://" + getHostName() + "/loadsettings", {method:'post'})
  .then(resp=>{
    if (!resp.ok) throw Error(rersp.statusText)
    return resp.json()
  })
  .then(data=>{
    if (data==null) return
 
    generalsettings.value=JSON.parse(data.generalsettings)
    managerauthorisations.value =JSON.parse(data.managersettings)
    printingprofiles.value=JSON.parse(data.printingsettings)
    
    securitysettings.value=data.securitysettings!= "" ? JSON.parse(data.securitysettings) : {disconnectmode:'two'}
   
    
   
  })
  .catch(ex=>{
    console.error(ex)
  })
}
function showComment(e){
  showcommentdialog.value=true
  selecteditem.value=e
  comment.value=e.comment
  
}
function saveComment(){
  selecteditem.value.comment=comment.value
 
}
async function saveSecurityOptions(){
  const form = new FormData()
  form.append("securitysettings", JSON.stringify(securitysettings.value))
  await fetch("https://" + getHostName() + "/savesecuritysettings", {method:'post', body: form})
  .then(resp=>{
    if (!resp.ok) throw Error(rersp.statusText)
    return resp.json()
  })
  .then(data=>{
    if (data==null) return
    alertHandlerTwo(data)
  })
  .catch(ex=>{
    console.error(ex)
  })
}
async function saveGeneralInformations(){
  
  const form = new FormData()
  form.append("generalsettings", JSON.stringify(generalsettings.value))
  await fetch("https://" + getHostName() + "/savegeneralsettings", {method:'post', body:form})
  .then(resp=>{
    if (!resp.ok) throw Error(rersp.statusText)
    return resp.json()
  })
  .then(data=>{
    if (data==null) return
     alertHandlerTwo(data)
  })
  .catch(ex=>{
    console.error(ex)
  })
}
function emptyTicketTwo(){
  subscriberInfos.value = null
  store.Member=null
}

function newTicketTwo(){
  subscriberInfos.value = null
  store.Member=null
  if (loggeduser.value.Role=="Gérant")
    selectedwaiter.value=''
  
  const newticketbutton= document.getElementById("newticketbutton")
 
  newticketbutton.blur()
 
}



const getTodayDate = () => {
      const today = new Date();
      const year = today.getFullYear();
      const month = String(today.getMonth() + 1).padStart(2, '0'); 
      const day = String(today.getDate()).padStart(2, '0');
      return `${year}-${month}-${day}`;
};
function formatDateToDDMMYYYYHHMMSS(dateString) {
  const date = new Date(dateString);

  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0'); 
  const year = date.getFullYear();

  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  return `${day}/${month}/${year} ${hours}:${minutes}:${seconds}`;
}
function isElectron() {
  return navigator.userAgent.toLowerCase().includes("electron")
}



function isRunningFromGUIApp(){
  if ((window.chrome && window.chrome.webview) || (window.nadilpos || isElectron()) )
    return true
  return false
}
onMounted(async () => {
  if (!isRunningFromGUIApp() ){
   
   window.location.href="notauthorised"
   return
  }
  document.documentElement.style.overflowY = 'hidden';
  document.body.style.overflowY = 'hidden';
  
  theme.value = localStorage.getItem('theme') || 'dark'
  
  
  await getUsers();
 
  const otpfirstchar= document.querySelector('[aria-label="Please enter OTP character 1"]')
  if(otpfirstchar!= null)
    otpfirstchar.focus()
  startdatemodel.value= getTodayDate()
  starttimemodel.value="00:00:00"
  enddatemodel.value=getTodayDate()
  ordersstatus.value='TOUT'
  initSortable();
});
const confirmstartendsessiondialog =ref(false)
const confirmstartendsessiontext =ref("")
async function confirmStartOrEndSession(){
  const form = new FormData()
  
  if (loggeduser.value.Username==null) return
  form.append("username",loggeduser.value.Username)
  await fetch("https://" + getHostName() + "/changesessionstatus",{method:'post', body:form})
  .then(r=>{
    if(!r.ok) throw Error(r.statusText)
      return r.json()
  })
  .then(async d=>{
    
      await getSessionsStatus()

  })
  .catch(ex=>{console.error(ex)})
}
async function showSituation(){
  
  showsituation.value=!showsituation.value
  usersnames.value=[]
  if (loggeduser.value?.Role== "Serveur"){
    
    usersnames.value.push(loggeduser.value?.Username)
    selectedperson.value= loggeduser.value?.Username
    userSelectionChanged(loggeduser.value?.Username)
    return
  }
  users.value?.forEach((w) => {
        usersnames.value.push(w.Username)
       
  });
  
 
}
function getDateTime(datetime){
 
  if (datetime==undefined) return ""
  if (datetime=="0001-01-01T00:00:00Z" || datetime=="" ){
    return ""
  }
  const date = new Date(datetime);
  const formattedDate = date.toISOString().split('T')[0]
  const formattedTime = date.toISOString().split('T')[1].slice(0, 5)
  return formattedDate + " " + formattedTime
}
const situation = ref("")
const situationToPrint = ref({})
async function userSelectionChanged(user){
  situationToPrint.value={}
  if (user=="") return
  
  if (sessions.get(user)?.start== undefined)  return
  const date = new Date(sessions.get(user).start);
  const formattedDate = date.toISOString().split('T')[0]
  const formattedTime = date.toISOString().split('T')[1].slice(0, 5)
  const form = new FormData()
  
  form.append("start", formattedDate + " " + formattedTime )
  form.append("username", user)
  await fetch("https://" + getHostName() + "/getsituation", {method:'post', body:form})
  .then(r=>{
    if (!r.ok) throw Error(r.statusText)
    return r.json()
  })
  .then(data=>{
    if (data== null) return
    situation.value=data
    
    situationToPrint.value={username:user,
                                totalpaid:situation.value?.totalpaid,
                                countpaid:situation.value?.countpaid,
                                totalinpaid:situation.value?.totalinpaid,
                                countinpaid:situation.value?.countinpaid,
                                totalcanceled:situation.value?.totalcanceled,
                                countcanceled:situation.value?.countcanceled,
                                total:situation.value?.total,
                                sessionstart:sessions.get(user)?.start,
                                sessionend:sessions.get(user)?.end

                                }
    
    
  })
  .catch(ex=>{
    console.error(ex)
  })
  
}
function showVirtualkeyBoard(){
  if (window.chrome.webview!=null)
    window.chrome.webview.postMessage("vkeyboard")
 
}
function statuschanged(){
  filtredorders.value=[]
  if (ordersstatus?.value != "TOUT"){
      
      
      orders.value.forEach(order=>{
        if (order.status == ordersstatus.value){
          filtredorders.value.push(order)
        }

      })
      
      return
    }
    filtredorders.value=orders.value
}

async function authOff(){

  await informServer()
  authenticated.value=false
  selectedwaiter.value=''
  startdatemodel.value= getTodayDate()
  starttimemodel.value="00:00:00"
  enddatemodel.value=getTodayDate()
  ordersstatus.value='TOUT'
  sessionStorage.clear()
  const otpfirstchar= document.querySelector('[aria-label="Please enter OTP character 1"]')
  if(otpfirstchar!= null)
    otpfirstchar.focus()
}
const showcommentdialog = ref(false)
const showalerttwo =ref(false)
const alerttexttwo = ref('')
const alertcardclasstwo = ref('')
async function alertHandlerTwo(data){

      showalerttwo.value = true;
      switch (data.Prefix) {
        case "ERROR":
        case "EMPTY":
          alerttexttwo.value = data.Message + `  ${data.data}`;
          alertcardclasstwo.value = "text-center errorclass";

          break;
        case "INSERT":
        case "EDIT":
        case "DELETE":
        case "UPDATE":
        
          alerttexttwo.value = data.Message ;
          alertcardclasstwo.value = "text-center successclass";
          await getOrders()
         

          break;
        case "SUCCESS":
        alerttexttwo.value = data.Message ;
        alertcardclasstwo.value = "text-center successclass";
          break;
      }
      setTimeout(() => {
            showalerttwo.value=false
      }, 4000);
      
      
     
}
async function quitApp(){
  if (window.nadilpos) {
    window.nadilpos.showMessage("close");
    return
  }
  
  if (window.chrome.webview!=null){
   
    window.chrome.webview.postMessage("close");
    
  }
  await authOff()

  
}
async function minimizeApp(){
  if (window.chrome.webview!=null){
    window.chrome.webview.postMessage("minimize");
    
  }

  
}

async function getOrders(){
  filtredorders.value=[]
  orders.value=[]
  const form = new FormData()
  form.append("startdate", startdatemodel.value )
  form.append("enddate", enddatemodel.value )
  form.append("starttime", starttimemodel.value)
  form.append("endtime", endtimemodel.value)
  await fetch("https://" + getHostName() + "/getOrders", {method:'post', body:form})
  .then(r=>{
    if (!r.ok) throw Error(r.statusText)
    return r.json()
  })
  .then(data=>{
    if (data==null) return
    orders.value=data
    
    if (ordersstatus?.value != "TOUT"){
      
      
      orders.value.forEach(order=>{
        if (order.status == ordersstatus.value){
          order.creationdate=formatDateToDDMMYYYYHHMMSS(order.creationdate)
          filtredorders.value.push(order)
        }

      })
      
      return
    }
    orders.value.forEach(order=>{
       
      order.creationdate=formatDateToDDMMYYYYHHMMSS(order.creationdate)
      filtredorders.value.push(order)
       

    })
    console.log(orders.value)
    
    
    
   
    
    
  })
  .catch(ex=>{
    console.error(ex)
  })
}
async function showOrdersList(){ 

  await getOrders()
  showorderslist.value =! showorderslist.value
 
 
}
async function informServer(){
  await fetch ("https://" + getHostName() + "/endsession")
  
  .catch(ex=>{console.error(ex)})
}
window.addEventListener('beforeunload', () => {
    document.cookie.split(";").forEach((cookie) => {
        document.cookie = cookie.replace(
            /^ +/,
            ""
        ).replace(/=.*/, `=;expires=${new Date(0).toUTCString()};path=/`);
    });
});

const sessionstatustext =ref('Commencer')
const sessionstatusicon =ref('mdi-play')
const sessionstatusclass=ref('statusplay')
const sessionsmap = reactive(new Map());
const sessions = reactive(new Map());
async function getSessionsStatus(){
  
  await fetch("https://" + getHostName() + "/getsessionsstatus",{method:'post'})
      .then(r=>{
        if(!r.ok) throw Error(r.statusText)
          return r.json()
      })
      .then(d=>{
        if (d== null) return
         
         d.forEach(session=>{
          sessions.set(session.username, session)
          
          sessionsmap.set(session.username, session.end =="0001-01-01T00:00:00Z" ? false:true)
          
          if(session.username == loggeduser.value?.Username){
            const date = new Date(session.start);
            const formattedDate = date.toISOString().split('T')[0]
            const formattedTime = date.toISOString().split('T')[1].slice(0, 5)
            startdatemodel.value= formattedDate
            starttimemodel.value=formattedTime
            if (session.end =="0001-01-01T00:00:00Z"){
              sessionstatustext.value="Finir"
              sessionstatusicon.value="mdi-stop"
              sessionstatusclass.value="statusstop"
              sessionended.value=true
              return
            }
           

            sessionstatustext.value="Commencer"
            sessionstatusicon.value="mdi-play"
            sessionstatusclass.value="statusplay"
            sessionended.value=false 
            return                                  
          }
          
           
          
              
              
         })
         
         if (sessionsmap.get(loggeduser.value?.Username)==undefined){
            
            sessionstatustext.value="Commencer"
            sessionstatusicon.value="mdi-play"
            sessionstatusclass.value="statusplay"
            sessionended.value=false 
            sessionsmap.set(loggeduser.value?.Username,true)
          }

      })
      .catch(ex=>{console.error(ex)})
      
}
async function changeSessionStatus(){
  confirmstartendsessiontext.value = sessionstatustext.value =="Commencer" ? "Voulez-vous continuer et démarrer une nouvelle session?" : "Êtes-vous sûr(e) de vouloir arrêter la session en cours?"
 confirmstartendsessiondialog.value=true

  

}
async function authOK(user){
  
  authenticated.value=true
  loggeduser.value=user
  await loadSettings()
  await getSessionsStatus()
  
  sessionStorage.setItem("user",JSON.stringify({username:user.Username, role:user.Role }))
  if (user.Role=='Serveur')
      selectedwaiter.value = user.Username
  
}

async function openSession(user){
      
      let form = new FormData()
      form.append("user", user)
      await fetch("https://" + getHostName() + "/opensession", {
        method:'post',
        body:form
      })
      .then(r=>{
        if (!r.ok) throw Error(r.statusText)
        return r.json()
        

      })
      .then(async data=>{
        await loadSettings()
        await getSessionsStatus()
      })
      .catch(ex=>{
        console.error(ex)
      })
}
document.addEventListener(
  "keydown",
  function (event) {
    
    keyssaved += event.key;
    lenghtcounter++;
    waiters.value.forEach((w) => {
      if (keyssaved.substring(0, lenghtcounter) === w.Code && authenticated.value==true && loggeduser.value?.Role=="Gérant" ) {
        if (w.Role=="Serveur")
          selectedwaiter.value = w.Username;
        
        keyssaved = "";
        lenghtcounter = 0;
        return;
      }
      if (keyssaved.substring(0, lenghtcounter) === w.Code && authenticated.value==false) {
        authenticated.value=true
        if (w.Role=="Serveur")
          selectedwaiter.value = w.Username;
        loggeduser.value=w
        keyssaved = "";
        lenghtcounter = 0;
        sessionStorage.setItem("user",JSON.stringify({username:w.Username, role:w.Role }))
        openSession(w.Username)
       
        return;
      }

    });

    if (event.key == "Enter") {
      keyssaved = "";
      lenghtcounter = 0;
    }
    setTimeout(() => {
      keyssaved = "";
      lenghtcounter = 0;
    }, 1000);
  },
  true
);
let rfidkeystrokes = ""
function enableRFIDReading() {
  document.addEventListener("keydown", handleRFIDInput);
 
}

function disableRFIDReading() {
  document.removeEventListener("keydown", handleRFIDInput);
  
  rfidkeystrokes = ""; 
}

function handleRFIDInput(event) {
  if (event.key.length == 1) {
    rfidkeystrokes+= event.key;
  }

  
  
  
  if (event.key == "Enter") {
    getSubscriberFromACM(rfidkeystrokes);
    rfidkeystrokes = "";
      
  }
  setTimeout(() => {
    rfidkeystrokes= "";
  }, 2000);
}
const getSubscriberFromACM = async (rfid) => {
  const form = new FormData();
  form.append("command", "getSubscriberInfosFromRFID");
  form.append("rfid_id", rfid);
  await fetch("https://" + window.location.hostname + ":8080/api", {
    method: "post",
    body: form,
  })
    .then((r) => {
      if (!r.ok) throw Error(r.statusText);
      return r.json();
    })
    .then((response) => {
      if (response == null) return;
      if (response.status =="SUCCESS"){
        if (response.data.hasValidSubscription == false){
          alertHandlerTwo({Prefix:"ERROR", Message: "Abonnement non valide", data:""})
          subscriberInfos.value = null;
          store.Member = null;
          return
        }
       
        alertHandlerTwo({Prefix:"SUCCESS", Message: "Abonnement valide", data:""})
        subscriberInfos.value = response.data;
        store.Member = response.data;
       
        
        
      }
    })
    .catch((e) => console.log(e));
};
async function getUsers() {
  waiters.value = [];
  waitersnames.value = [];
  usersnames.value=[];
  await fetch("https://" + getHostName() + "/getusers")
    .then((r) => {
      if (!r.ok) throw Error(r.statusText);
      return r.text();
    })
    .then((data) => {
      users.value = JSON.parse(data);
      
      users.value?.forEach((w) => {
        usersnames.value.push(w.Username)
        
          waiters.value.push({
            Username: w.Username,
            Role: w.Role,
            Code: w.Code,
          });
          waitersnames.value.push(w.Username);
        
      });
    })
    .catch((e) => console.log(e));
}
const confirmcancelorderdialg=ref(false)
const confirmcancelordertext=ref("")
const ordertocancel=ref(null)
function confirmCancelOrder(order){
  confirmcancelordertext.value="Continuer pour annuler la commande n° : " + order.number + " ?"
  ordertocancel.value=order
  confirmcancelorderdialg.value=true
}
async function cancelOrder() {
  if (ordertocancel.value?.status=="ANNULE"){
    alertHandlerTwo({Prefix:"ERROR", Message:"Ce ticket est déjà annulé.",data:"CANCELED_ORDER"})
    return
  }
  
   
   const form = new FormData() 
  
   form.append("number", ordertocancel.value?.number)
   
   await fetch ("https://" + getHostName() + "/cancelorder", {method:'post', body: form})
   .then(r=>{
          if (!r.ok) throw Error(r.status)
          return r.json()
        
    })
    .then(data=>{
      if (data==null) return
       alertHandlerTwo(data)
    })
}
const confirmcheckorderdialog =ref(false)
const confirmcheckordertext =ref("")
const ordertocheck= ref(false)
function confirmCheckOrder(order){
  confirmcheckordertext.value= "Voulez-vous valider le paiement de la commande  N° " +  order.number + " ?"
  confirmcheckorderdialog.value=true
  ordertocheck.value=order
}
async function checkConfirmation(){
  if (ordertocheck.value?.status=="ANNULE"){
    alertHandlerTwo({Prefix:"ERROR", Message:"Ce ticket a été annulé.",data:"CANCELED_ORDER"})
    return
  }
  if (ordertocheck.value?.status=="PAYE"){
    alertHandlerTwo({Prefix:"ERROR", Message:"Ce ticket est déjà payé.",data:"PAID_ORDER"})
    return
  }
 
   const form = new FormData() 
  
   form.append("number", ordertocheck.value?.number)
   form.append("creator", ordertocheck.value?.creator)
   form.append("waiter", ordertocheck.value?.waiter)
   await fetch ("https://" + getHostName() + "/checkoutorder", {method:'post', body: form})
   .then(r=>{
          if (!r.ok) throw Error(r.status)
          return r.json()
        
    })
    .then(data=>{
      if (data==null) return
       alertHandlerTwo(data)
    })
}

function showContent(item){
  showcontent.value=true
  
  selectedorder.value=item

}
function confirmDeleteUser(user) {
  confirmtext.value =
    "Êtes-vous sûr de vouloir supprimer l'utilisateur " + user.Username + "?";
  usertodelete.value = user;
  confirmdeleteuderdialog.value = !confirmdeleteuderdialog.value;
}
async function deleteUser() {
  const form = new FormData();
  
  form.append("id", usertodelete.value.ID);
  await fetch("https://" + getHostName() + "/deleteuser", {
    method: "post",
    body: form,
  })
    .then((r) => {
      if (!r.ok) throw Error(r.statusText);
      return r.json();
    })
    .then((data) => {
     
      inform();
      confirmdeleteuderdialog.value = false;
    })
    .catch((e) => console.log(e));
}

function openManageUsersDialog() {
  isManageUsersDialogOpen.value = !isManageUsersDialogOpen.value;
}
function closeManageUsersDialog(value) {
  isManageUsersDialogOpen.value = value;
}
function openEditUsersDialog() {
 
  isEditUsersDialogOpen.value = !isEditUsersDialogOpen.value;
}
function closeEditUsersDialog(value) {
  isEditUsersDialogOpen.value = value;
}
function openAddCategoryDialog() {
  addcatbutton = true;
  isCategoryDialogOpen.value = !isCategoryDialogOpen.value;
}
function closeCategoryDialog(value) {
  isCategoryDialogOpen.value = value;
}
function openAddProductDialog() {
  isProductDialogOpen.value = !isProductDialogOpen.value;
}
function closeProductDialog(value) {
  isProductDialogOpen.value = value;
}
function openEditProductDialog() {
  isEditProductDialogOpen.value = !isEditProductDialogOpen.value;
}
function closeEditProductDialog(value) {
  isEditProductDialogOpen.value = value;
}


window.addEventListener(
  "contextmenu",
  function (e) {
    e.preventDefault();
  },
  false
);
async function inform() {

  await getUsers();
}

function editCategory(cat) {
  editcat = cat;
  isCategoryDialogOpen.value = true;
}
function showMenu() {
  
  rail.value = !rail.value;
  drawer.value = !drawer.value;
}
async function printSituation(){
  if (situationToPrint.value?.sessionend == "0001-01-01T00:00:00Z" || situationToPrint.value?.sessionend == ""){
    alertHandlerTwo({Prefix:"ERROR",Message:"Votre session n'est pas encore finie. Veuillez mettre fin à la session en cours avant d'imprimer.",data:""})
    return
  }
  if (situationToPrint.value?.username == undefined){
    alertHandlerTwo({Prefix:"ERROR",Message:"Vous n'avez aucune session.",data:""})
    return
  }
  const form = new FormData()
  form.append("situation", JSON.stringify(situationToPrint.value))

  
  await fetch("https://" + getHostName() + "/printsituation", {method:'post', body:form})
  .then(r=>{
    if (!r.ok) throw Error(r.statusText)
    return r.json()
  })
  .then(data=>{
    if (data == null) return
    
  })
  .catch(ex=>{
    console.log(ex)
  })
}
function getMenuDisabled(btnvalidatedisabled, ticketvalidate,sessionend, isadmin,managerauthorisations, selectedwaiter){
  if (isadmin) return false
  if (sessionend) return true
  if (managerauthorisations.managercancreatetickets ==false && selectedwaiter =="") return true
  if (btnvalidatedisabled && !ticketvalidate) return true
  



 

}
</script>
<script>
import { getHostName, filtredMenu, initSortable, getSortedFiltredMenu, onUpdatefocused, store} from "./Utils";
import OrderContentComponent from "./components/OrderContentComponent.vue";
import ManageUsersDialogComponent from "./components/ManageUsersDialogComponent.vue";
import EditUsersDialogComponent from "./components/EditUsersDialogComponent.vue";
import AuthNumPadComponent from "./components/AuthNumPadComponent.vue"
import NewCategoryDialogComponent from "./components/NewCategoryDialogComponent.vue";
import NewProductDialogComponent from "./components/NewProductDialogComponent.vue";
import EditProductDialogComponent from "./components/EditProductDialogComponent.vue";


let oldtab = 0;

const axios = require("axios").default;


export default {
  components: { ManageUsersDialogComponent },

  mounted: async function () {
    
    await this.getCategories();
    
   
  },
  data() {
    return {
    

      deleteproductconfirmdialog: false,
      overlay:true,
      authenticated_:true,
      producttodelete: "",
      ticketvalidate: true,
      showvkeyboard: true,
      editelement: null,
      showpad: false,
      actiononproduct: false,
      btnvalidatedisabled: true,
      btnvalidate: false,
      newcategoryname: "",
      showadminprompt: false,
      showalert: false,
      confirmdialog: false,
      adminclass: "adminclasslogout",
      adminloginicon: "mdi-lock",
      sessionicon:"mdi-play",
      confirmationtext: "",
      selectedcategory: "",
      alertcardclass: "text-center",
      alerttext: "",
      alerttitle: "",
      adminpassword: "",
      ticketnumber: "",
      isadmin: false,
      paramsdialog: false,
      productdialog: false,
      sheet:false,
      tab: -1,
      categorytodelete: "",
      input: "",
      categorydialog: false,
      selectedcategoryID: -1,
      currentelement: null,
      count: 0,
     
   
      categories: [],
      menu: [],
     
      ticket: [
        { content: [], number: "", totalticket: 0, waiter: '', member: '', memberid: '' },
      ],
    };
  },
  computed: {
    store() {
      return store // make it available in template and watchers
    }
  },
  watch: {
    tab: async function (val) {
      if (val == 0) {
        if (this.menu.length == 0) {
          await this.getMenu();
        }
        this.showFavorites();
        return;
      }

      this.showElements(this.selectedcategory.categoryid);
    
    },
   
   'store.Member': {
      handler(val) {
        if (val != null) {
          console.log(val)
          this.ticket[0].member = val?.name;
          this.ticket[0].memberid = val?.subscriberId;
          this.ticket[0].content.forEach(c => {
          
            c.price = this.selectPrice(this.menu.find(m => m.productid == c.id))
          })
          this.updateTotalTicket()
          console.log(this.ticket[0])
        }
      },
      immediate: true,
      deep: true
    },
    ticket: {
      handler: function (val) {
        this.btnvalidatedisabled = !val[0].content.length > 0;
        this.ticketvalidate = true;
      },
      deep: true,
    },
  },

  methods: {
   
    logout(){
      
      this.ticket[0].content=[];
      this.ticket[0].totalticket =0;
      this.ticket[0].waiter=''
      this.ticket[0].number=''
      this.ticketnumber = '';
      this.isadmin=false
      this.adminloginicon = "mdi-lock"
      this.adminclass='adminclasslogout'
      
      
    },
    async adminLogin() {
      const form = new FormData(); 

      form.append("pwd", this.adminpassword);
      await fetch("https://" + getHostName() + "/adminlogin", {
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
          this.showadminprompt = false;
          if (data.data == "ACCES REFUSE") {
            this.alerttext = data.Message;

            this.showalert = true;
            this.alertcardclass = "text-center errorclass";
            return;
          }
          this.adminclass = "adminclasslogin";
          this.adminloginicon = "mdi-lock-open";
          this.isadmin = true;
          this.adminpassword = "";
        })
        .catch((e) => {
          console.log(e);
        });
    },
    newTicket() {
      this.ticket[0].content = [];
      this.ticket[0].totalticket = 0;
      this.ticketvalidate = false;
      this.ticketnumber = "";
      
      document.body.focus()
    },
    emptyTicket() {
      this.ticket[0].content = [];
      this.ticket[0].totalticket = 0;
    },
    async validateTicket(sessions, ss, authoff) {
      
      const inputwaiter= this.$refs.waiterbox.$el.querySelector('[inputmode="none"]')
      if (inputwaiter.value!=''){
        
        
        if (sessions.get(inputwaiter.value)){
          this.alertHandler({Prefix:"ERROR", Message:"La session du serveur selectionné n'est pas démarrée.", data:""})
          return
        }
      }
      
      const hostname = getHostName();
      const form = new FormData();
      const creator = JSON.parse(sessionStorage.getItem('user')).username
      
      
      this.ticket[0].waiter=inputwaiter.value
      form.append("totalticket", this.ticket[0].totalticket);
      form.append("waiter", this.ticket[0].waiter);
      form.append("creator", creator)
      form.append("content", JSON.stringify(this.ticket[0].content));
      form.append("member", this.ticket[0].member);
      form.append("memberid", this.ticket[0].memberid);

      await fetch("https://" + hostname + "/saveTicket", {
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
          this.ticketvalidate = false;
          this.btnvalidatedisabled = true;

          this.ticketnumber = data.data.padStart(9, 0);
          if (ss.disconnectmode=="one"){
            authoff()
            this.logout()
            
          }
          
        })
        .catch((e) => {
          console.log(e);
        });
    },
    updateTotalTicket() {
      this.ticket[0].totalticket = 0;
      this.ticket[0].content.forEach((c) => {
        this.ticket[0].totalticket += c.count * this.selectPrice(this.menu.find(m => m.productid == c.id));
      });
    },
    validateInput() {
      const display = document.getElementById("display").innerHTML;
      if (display !== "Entrez un nombre" && display !== "") {
        if (display == "0") {
          this.removeFromTicket(this.currentelement.id);
          this.showpad = false;
          return;
        }
        this.currentelement.count = parseInt(display);
        this.updateTotalTicket();
        this.showpad = false;
      } else {
        alert("Nombre incorrecte");
      }
    },
    pressButton(value) {
      const display = document.getElementById("display");
      if (display.innerHTML === "Entrez un nombre") {
        display.innerHTML = value;
      } else {
        display.innerHTML += value;
      }
    },

    deleteLast() {
      const display = document.getElementById("display");
      display.innerHTML = display.innerHTML.slice(0, -1) || "Entrez un nombre";
    },

    clearDisplay() {
      document.getElementById("display").innerHTML = "Entrez un nombre";
    },
    showPad(e) {
      this.showpad = true;
      this.currentelement = e;
    },
    async refresh() {
      await this.getCategories();
      await this.getMenu();
    },
    getImageURL(i) {
      return "https://" + getHostName() + "/img/" + i;
    },
 
    async printReceipt(ss, authoff){
      const form = new FormData()
      form.append("number", this.ticketnumber)

      await fetch("https://"+ getHostName()+ "/printreceipt", {method:"post",body:form})
      .then(r=>{
        if (!r.ok) throw Error(r.statusText)
        return r.json()
      })
      .then(data=>{
        
         if (ss.disconnectmode=="two"){
            authoff()
            this.logout()
            
        }
      })
      .catch(err=>{
        console.log(err)
      })
    },
    getMenu: async function () {
      this.menu = [];
      await fetch("https://" + getHostName() + "/getMenu")
        .then((r) => {
          if (!r.ok) throw Error(r.statusText);
          return r.json();
        })
        .then((data) => {
          this.menu = data; 
          
         
        })
        .catch((e) => console.log(e));
    },
    getCategories: async function () {
      this.categories = [];
      await axios
        .post("https://" + getHostName() + "/getcategories", {
          withCredentials: true,
        })
        .then((r) => {
          this.categories = r.data;
        })
        .catch((e) => console.log(e));
    },
    selectPrice: function (el) {
      
      if (store.Member == null) {
        return parseFloat(el?.price);
        
      }
      if (store.Member?.isGold && store.Member?.remainingConsumables>0){
        
        return el?.goldsubscriberprice ? parseFloat(el?.goldsubscriberprice): parseFloat(el?.price);
       
        
     
       
      }
      return el?.subscriberprice ? parseFloat(el?.subscriberprice): parseFloat(el?.price);
    },

    elementClick: function (el) {
      
      
      if (this.isadmin) {
        this.actiononproduct = true;

        return;
      }
      let found = false;
      this.ticket[0].totalticket =
        this.ticket[0].totalticket + this.selectPrice(el);
      
      this.ticket[0].content.forEach((e) => {
        if (e.id == el.productid) {
          e.count++;
          found = true;
          return;
        }
      });
      if (found == false)
        this.ticket[0].content.push({
          id: el.productid,
          count: 1,
          text: el.text,
          price: this.selectPrice(el),
          categoryid : el.category,
          comment: "",
        });

      var listticket = document.getElementById("listticket");
      listticket.lastElementChild.scrollIntoView(false);
    },
    
    async alertHandler(data) {
      this.showalert = true;
      switch (data.Prefix) {
        case "ERROR":
        case "EMPTY":
          this.alerttext = data.Message + `  ${data.data}`;
          this.alertcardclass = "text-center errorclass";

          break;
        case "INSERT":
        case "EDIT":
        case "DELETE":
        case "UPDATE":
          this.alerttext = data.Message + ` : ${data.data}`;
          this.alertcardclass = "text-center successclass";

          if (data.Message == "Le ticket a été enregistré")
            this.ticket[0].number = data.data;

          break;
      }
      await this.getCategories();
    
      await this.getMenu();
      if (oldtab == 0 && (data.Message != "La catrégorie a été enregistrée" || data.Message!="La catrégorie a été modifiée") ) {
        this.showFavorites();
        return;
      }
      this.tab = oldtab;
      this.showElements(this.selectedcategory.categoryid);
      setTimeout(() => {
            this.showalert=false
      }, 3000);
    },
    async showElements(catID) {
      if (catID != "") {
        this.selectedcategoryID = catID;
      }
      
      if (this.menu.length == 0)  
        await this.getMenu();


      filtredMenu.value = [];
      this.menu.forEach((e) => {
        if (e.category == catID) {
          filtredMenu.value.push(e);
        }
      });
      initSortable(catID);
      if (localStorage.getItem("filtredMenu" + catID)){
        filtredMenu.value = getSortedFiltredMenu(catID);
        return
      }
      localStorage.setItem("filtredMenu" + catID, JSON.stringify(filtredMenu.value));
    },
    async showFavorites() {
      filtredMenu.value = [];
     
     
      
      this.menu.forEach((e) => {
        if (e.isfavorite == true) {
          filtredMenu.value.push(e);
        }
      });
      initSortable(-1);
      if (localStorage.getItem("filtredMenu-1")){
        filtredMenu.value = getSortedFiltredMenu(-1);
        return
      }
      localStorage.setItem("filtredMenu-1", JSON.stringify(filtredMenu.value));
    },
    removeFromTicket: function (id) {
      this.ticket[0].content.forEach((e, i) => {
        if (e.id == id) {
          this.ticket[0].content.splice(i, 1);
          this.ticket[0].totalticket -= parseFloat(e.price) * e.count;
          return;
        }
      });
    },
    increaseCount: function (e) {
      e.count++;
      this.ticket[0].totalticket += parseFloat(e.price);
    },
    reduceCount: function (e) {
      if (e.count == 1) {
        e.count--;
        this.ticket[0].totalticket -= parseFloat(e.price);
        this.removeFromTicket(e.id);
        return;
      }

      e.count--;
      this.ticket[0].totalticket -= parseFloat(e.price);
    },
    showTicketOnLowResolution() {
      this.sheet = true;
    },
    async showadminlogin(func) {
      if (this.adminloginicon == "mdi-lock") {
        this.showadminprompt = true;
        setTimeout(() => {
          this.$refs.adminpass.focus();
        }, 200);
        return;
      }
      this.isadmin = false;
      this.adminclass = "adminclasslogout";
      this.adminloginicon = "mdi-lock";
      func()
    },
    showAddNewProductDialog() {
      oldtab = this.tab;

      this.productdialog = true;
    },

    async confirmDelete() {
      await axios
        .post("https://" + getHostName() + "/deletecategory", {
          categoryid: this.categorytodelete.categoryid,
        })
        .then(async (r) => {
          this.alertHandler(r.data);
        })
        .catch((e) => console.log(e));
    },
    async confirmDeleteProduct() {
      oldtab = this.tab;
      await axios
        .post("https://" + getHostName() + "/deleteProduct", {
          productid: this.producttodelete.productid,
        })
        .then(async (r) => {
          this.alertHandler(r.data);

          await this.getMenu();
        })
        .catch((e) => console.log(e));
    },

    deleteCategory(cat) {
      this.categorytodelete = cat;
      this.confirmdialog = true;
      this.confirmationtext = `Etes-vous sûr de vouloir supprimer la catégorie ${cat.categoryname}? `;
    },
    deleteProduct(element) {
      this.producttodelete = element;
      this.deleteproductconfirmdialog = true;
      this.confirmationtext = `Etes-vous sûr de vouloir supprimer le produit ${element.text}? `;
    },

  },
};
</script>
<style>
.situation {
    display: flex;
    color: black;
    justify-content: center; 
    align-items: center;    
    width: 400px;
    
}
.v-main{
  --v-layout-bottom : 0px;
}
.errorclass {
  background-color: chocolate;
}
.successclass {
  background-color: darkseagreen;
}
.v-list-subheader__text {
  margin-top: 5px;
  margin-left: 5px;
}
.v-list-subheader {
  display: grid;
}
.adminclasslogin {
  color: darkgreen;
}
.adminclasslogout {
  color: darkred;
}
.statusplay{
  color:rgb(50, 112, 31);
}
.statusstop{
  color:rgb(112, 31, 31);
}
.v-card-actions {
  display: inline !important;
}
.v-col-md-2 {
  flex: 0 0 16.6666666667%;
  max-width: 50%;
}
.v-chip.v-chip--size-small {
  padding: 0 16px;
}
.no-select {
  -webkit-user-select: none;
  -ms-user-select: none;
  user-select: none;
}
.IMPAYE {
  background-color:orangered;
  text-align: center;
}
.PAYE {
  background-color: darkolivegreen;
  text-align: center;
}
.ANNULE {
 background-color: #888;
text-align: center;
}
.menu {
  width: 65%;
}
.showinlowresolution {
  visibility: hidden;
}
.centerstuff {
  position: absolute;

  top: 40%;
  left: 50%;
  transform: translate(-50%, -50%);
}
.card-element {
  height: 120px;
  width: 80px;
}
.card-element-no-drag{
  height: 120px;
  width: 80px;
  
}
.card-element-image {
  height: 90px;
  width: 80px;
  object-fit: cover;
}
.adminlogin {
  width: 145px;
  padding: 10px;
  font-size: 14px;
  text-align: center;
  border: 2px solid #ccc;
  border-radius: 5px;
  margin-bottom: 20px;
  background-color: #ffffff;
}



.display {
  width: 145px;
  padding: 10px;
  font-size: 14px;
  text-align: center;
  border: 2px solid #ccc;
  border-radius: 5px;
  margin-bottom: 20px;
  color: black;
  background-color: #ffffff;
}
.keypad {
  display: grid;
  grid-template-columns: repeat(3, 40px);
  gap: 10px;
}
.keypad button {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  font-size: 20px;
  background-color: #1c1d1d;
  color: rgb(243, 228, 228);
  border: none;
  cursor: pointer;
  outline: none;
  transition: background-color 0.2s;
}
.keypad button:active {
  background-color: #0056b3;
}
.keypad .control {
  background-color: #dc3545;
}
.keypad .control:active {
  background-color: #c82333;
}
.keypad .validate {
  background-color: #28a745;
  grid-column: 2 / span 1;
}
.keypad .validate:active {
  background-color: #218838;
}
.v-overlay__scrim {
  background: transparent;
}
.keyboarddialog {
  max-width: 750px;
}
@media only screen and (max-width:1900px){
  .v-col-md-2 {
    flex: 0 0 16.6666666667%;
    max-width: 20% !important;
}

.v-container {
        max-width: 100%;
    }
}

@media only screen and (max-width: 757px) {
  .v-col-md-2 {
    flex: 0 0 16.6666666667%;
    max-width: 40% !important;

}
.v-bottom-navigation .v-bottom-navigation__content>.v-btn {
  width:50px;
}
body {
  margin-right:5px;
  margin-left:5px;
}
  
.bottomnavbutton span {
  font-size: 9px;
  
}
  .NOMOBILE{
    display:none;
  }
  .menu {
    width: 100%;
  }
  .situation {
   width: 100vw;
   left:0%;
   max-width: 100%;
  }
 
  .keyboarddialog {
    display: none;
  }
  .card-element {
    height: 120px;
    width: 82px;
  }
  .card-element-image {
    height: 90px;
    width: 82px;
    object-fit: cover;
  }
  .ticket {
    visibility: hidden;
    display: none;
  }
  .showinlowresolution {
    visibility: visible;
  }
}
input[type="date"] {
      background-color: #1e1e1e;
      color: #ffffff;
      border: 1px solid #333;
      border-radius: 5px;
      padding: 10px;
      font-size: 12px;
      outline: none;
}
select {
  background-color: #333;
  color: white;
  border: 1px solid #555;
  border-radius: 4px;
  padding: 8px;
  font-size: 12px;
}

select:focus {
  outline: none;
  border-color: #888;
}

option {
  background-color: #333;
  color: white;
}

option:hover {
  background-color: #555;
  
}
.v-selection-control .v-label{
  font-size: 12px !important;
}
.v-col-md-2 {
  flex: 0 0 16.6666666667%;
  max-width: 10%;

}
.d-flex {
  display: flex;
  flex-wrap: wrap;
  gap: 16px; /* Ensures space between cards */
}

.card-element {
  cursor: grab;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

/* Dragging effect */
.card-element:active {
  cursor: grabbing;
}

.sortable-ghost {
  opacity: 0.4;
}

.card-element:hover {
  transform: scale(1.05);
}
.highlighted {
  color: chocolate;
}
.active{
  color: darkseagreen;
}
.inactive{
  color: darkred;
}
</style>
``` 
