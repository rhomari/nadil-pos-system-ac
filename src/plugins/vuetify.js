// Styles
import "@mdi/font/css/materialdesignicons.css";
import "vuetify/styles";

// Vuetify
import { createVuetify } from "vuetify";
import { VNumberInput } from 'vuetify/labs/VNumberInput'
export default createVuetify({
    components: {
      VNumberInput,
    },
    theme: {
      defaultTheme: localStorage.getItem('theme') || 'dark',
      themes: {
        bluegrey: {
          colors: {
            primary: '#4CAF93',
            secondary: '#FFB74D',
            accent: '#82B1FF',
            background: '#5d6d7e',
            surface: '#34495e',
            error: '#E57373',
            success: '#81C784',
            warning: '#FFD54F',
            info: '#64B5F6',
            textPrimary: '#212121',
            textSecondary: '#757575',
          },
          variables: {
            borderRadius: '12px',
            customShadow: '0px 4px 10px rgba(0, 0, 0, 0.1)',
          },
        },
        brown: {
          colors: {
            primary: '#3E2723',
            secondary: '#FFB74D',
            accent: '#82B1FF',
            background: '#8D6E63',
            surface: '#A1887F',
            error: '#F44336',
            success: '#66BB6A',
            warning: '#FFC107',
            info: '#29B6F6',
            textPrimary: '#212121',
            textSecondary: '#757575',
          },
          variables: {
            borderRadius: '10px',
            customShadow: '0px 4px 12px rgba(0, 0, 0, 0.1)',
          },
        },
        teal: {
          colors: {
            primary: '#004D40',
            secondary: '#FFB74D',
            accent: '#82B1FF',
            background: '#26A69A',
            surface: '#80CBC4',
            error: '#F44336',
            success: '#66BB6A',
            warning: '#FFC107',
            info: '#29B6F6',
            textPrimary: '#212121',
            textSecondary: '#757575',
          },
          variables: {
            borderRadius: '10px',
            customShadow: '0px 4px 12px rgba(0, 0, 0, 0.1)',
          },
        },
        purple: {
          colors: {
            primary: '#673AB7',
            secondary: '#E91E63',
            accent: '#FF4081',
            background: '#EDE7F6',
            surface: '#D1C4E9',
            error: '#D32F2F',
            success: '#4CAF50',
            warning: '#FFEB3B',
            info: '#3F51B5',
            textPrimary: '#1A237E',
            textSecondary: '#4A148C',
          },
          variables: {
            borderRadius: '8px',
            customShadow: '0px 5px 15px rgba(0, 0, 0, 0.2)',
          },
        },
        orange: {
          colors: {
            primary: '#FF5722',
            secondary: '#FFC107',
            accent: '#FF9800',
            background: '#FFCCBC',
            surface: '#FFAB91',
            error: '#D84315',
            success: '#2E7D32',
            warning: '#FF8F00',
            info: '#FF6F00',
            textPrimary: '#BF360C',
            textSecondary: '#5D4037',
          },
          variables: {
            borderRadius: '10px',
            customShadow: '0px 4px 14px rgba(0, 0, 0, 0.15)',
          },
        },
        deepblue: {
          colors: {
            primary: '#0D47A1',
            secondary: '#1565C0',
            accent: '#42A5F5',
            background: '#1E88E5',
            surface: '#64B5F6',
            error: '#B71C1C',
            success: '#2E7D32',
            warning: '#F57C00',
            info: '#1976D2',
            textPrimary: '#E3F2FD',
            textSecondary: '#BBDEFB',
          },
          variables: {
            borderRadius: '10px',
            customShadow: '0px 6px 16px rgba(0, 0, 0, 0.2)',
          },
        },
      },
    },
    
  })