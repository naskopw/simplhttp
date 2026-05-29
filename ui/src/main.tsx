import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { MantineProvider, createTheme } from '@mantine/core'
import { Notifications } from '@mantine/notifications'
import '@mantine/core/styles.css'
import '@mantine/notifications/styles.css'
import '@mantine/dropzone/styles.css'
import './index.css'
import App from './App.tsx'

const theme = createTheme({
  primaryColor: 'violet',
  fontFamily: 'var(--sans)',
  headings: {
    fontFamily: 'var(--heading)',
  },
})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <MantineProvider theme={theme} defaultColorScheme="auto">
      <Notifications />
      <App />
    </MantineProvider>
  </StrictMode>,
)
