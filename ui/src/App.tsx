import { useState, useEffect } from 'react'
import axios from 'axios'
import logo from './assets/logo.svg'
import {
  AppShell,
  Title,
  Text,
  Container,
  Group,
  Table,
  Button,
  Breadcrumbs,
  Anchor,
  ActionIcon,
  Stack,
  Notification,
  rem,
  Image,
  useMantineColorScheme,
  useComputedColorScheme,
  Paper,
  Box,
  Loader,
  Progress,
} from '@mantine/core'
import { Folder, File, Download, Upload, ArrowLeft, Sun, Moon } from 'lucide-react'
import { Dropzone } from '@mantine/dropzone'
import { useMediaQuery } from '@mantine/hooks'

interface FileInfo {
  name: string
  size: number
  isDir: boolean
  modTime: string
}

interface ServerConfig {
  maxSizeBytes: number
  readOnly: boolean
}

function formatBytes(bytes: number, decimals = 2) {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

function App() {
  const [currentPath, setCurrentPath] = useState('')
  const [files, setFiles] = useState<FileInfo[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [uploading, setUploading] = useState(false)
  const [uploadProgress, setUploadProgress] = useState(0)
  const [config, setConfig] = useState<ServerConfig>({ maxSizeBytes: 100 * 1024 ** 2, readOnly: false })
  const isMobile = useMediaQuery('(max-width: 768px)')
  const { setColorScheme } = useMantineColorScheme()
  const computedColorScheme = useComputedColorScheme('light', { getInitialValueInEffect: true })

  const fetchFiles = async (path: string) => {
    setLoading(true)
    setError(null)
    try {
      const response = await fetch(`/api/fs/${path}`)
      if (!response.ok) throw new Error('Failed to fetch files')
      const data = await response.json()
      setFiles(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  const fetchConfig = async () => {
    try {
      const response = await fetch('/api/config')
      if (response.ok) {
        const data = await response.json()
        setConfig(data)
      }
    } catch (err) {
      console.error('Failed to fetch config:', err)
    }
  }

  useEffect(() => {
    fetchConfig()
  }, [])

  useEffect(() => {
    fetchFiles(currentPath)
  }, [currentPath])

  const navigateTo = (name: string) => {
    const newPath = currentPath ? `${currentPath}/${name}` : name
    setCurrentPath(newPath)
  }

  const navigateBack = () => {
    const parts = currentPath.split('/')
    parts.pop()
    setCurrentPath(parts.join('/'))
  }

  const getFileUrl = (name: string) => {
    const path = currentPath ? `${currentPath}/${name}` : name
    return `/api/fs/${path}`
  }

  const handleDownload = (name: string) => {
    window.open(getFileUrl(name), '_blank')
  }

  const onDrop = async (files: File[]) => {
    setUploading(true)
    setUploadProgress(0)
    const formData = new FormData()
    files.forEach((file) => {
      formData.append('files', file)
    })

    try {
      await axios.post(`/api/upload/${currentPath}`, formData, {
        onUploadProgress: (progressEvent) => {
          if (progressEvent.total) {
            const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total)
            setUploadProgress(percentCompleted)
          }
        },
      })
      fetchFiles(currentPath)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed')
    } finally {
      setUploading(false)
      setUploadProgress(0)
    }
  }

  const breadcrumbs = currentPath.split('/').filter(Boolean)

  const toggleColorScheme = () => {
    setColorScheme(computedColorScheme === 'dark' ? 'light' : 'dark')
  }

  const rows = files.map((file) => (
    <Table.Tr key={file.name}>
      <Table.Td>
        <Group gap="xs">
          {file.isDir ? (
            <Folder size={18} color="var(--mantine-color-violet-filled)" />
          ) : (
            <File size={18} color="var(--mantine-color-gray-6)" />
          )}
          {file.isDir ? (
            <Anchor onClick={() => navigateTo(file.name)} fw={500}>
              {file.name}
            </Anchor>
          ) : (
            <Text size="sm">{file.name}</Text>
          )}
        </Group>
      </Table.Td>
      {!isMobile && (
        <>
          <Table.Td>
            <Text size="sm" c="dimmed">
              {file.isDir ? '-' : formatBytes(file.size)}
            </Text>
          </Table.Td>
          <Table.Td>
            <Text size="sm" c="dimmed">
              {new Date(file.modTime).toLocaleString()}
            </Text>
          </Table.Td>
        </>
      )}
      <Table.Td>
        {!file.isDir && (
          <ActionIcon
            variant="subtle"
            color="violet"
            onClick={() => handleDownload(file.name)}
            title="Download"
          >
            <Download size={16} />
          </ActionIcon>
        )}
      </Table.Td>
    </Table.Tr>
  ))

  const mobileCards = files.map((file) => (
    <Paper key={file.name} withBorder p="md" radius="md">
      <Group justify="space-between" align="start">
        <Group gap="sm" style={{ flex: 1 }}>
          {file.isDir ? (
            <Folder size={24} color="var(--mantine-color-violet-filled)" />
          ) : (
            <File size={24} color="var(--mantine-color-gray-6)" />
          )}
          <Box style={{ flex: 1 }}>
            {file.isDir ? (
              <Anchor onClick={() => navigateTo(file.name)} fw={600} size="md">
                {file.name}
              </Anchor>
            ) : (
              <Text fw={500} size="md">
                {file.name}
              </Text>
            )}
            <Group gap="xs" mt={4}>
              <Text size="xs" c="dimmed">
                {file.isDir ? 'Directory' : formatBytes(file.size)}
              </Text>
              <Text size="xs" c="dimmed">•</Text>
              <Text size="xs" c="dimmed">
                {new Date(file.modTime).toLocaleDateString()}
              </Text>
            </Group>
          </Box>
        </Group>
        {!file.isDir && (
          <ActionIcon
            variant="light"
            color="violet"
            onClick={() => handleDownload(file.name)}
            size="lg"
          >
            <Download size={20} />
          </ActionIcon>
        )}
      </Group>
    </Paper>
  ))

  return (
    <AppShell header={{ height: 60 }} padding="md">
      <AppShell.Header>
        <Container size="xl" h="100%">
          <Group h="100%" px="md" justify="space-between">
            <Group>
              <Image src={logo} h={30} w={30} alt="simplhttp logo" />
              <Title order={3} visibleFrom="xs">simplhttp</Title>
            </Group>
            <Group>
              {loading && <Loader size="sm" />}
              <ActionIcon
                onClick={toggleColorScheme}
                variant="default"
                size="lg"
                aria-label="Toggle color scheme"
              >
                {computedColorScheme === 'dark' ? <Sun size={18} /> : <Moon size={18} />}
              </ActionIcon>
            </Group>
          </Group>
        </Container>
      </AppShell.Header>

      <AppShell.Main>
        <Container size="xl">
          <Stack gap="lg">
            <Paper withBorder p="sm" radius="md" bg="var(--mantine-color-body)">
              <Group justify="space-between">
                <Breadcrumbs
                  styles={{
                    breadcrumb: { fontSize: isMobile ? rem(14) : rem(16) },
                    separator: { marginInline: rem(8) },
                  }}
                >
                  <Anchor onClick={() => setCurrentPath('')} fw={currentPath === '' ? 700 : 400}>
                    Root
                  </Anchor>
                  {breadcrumbs.map((item, index) => {
                    const path = breadcrumbs.slice(0, index + 1).join('/')
                    return (
                      <Anchor
                        key={path}
                        onClick={() => setCurrentPath(path)}
                        fw={index === breadcrumbs.length - 1 ? 700 : 400}
                      >
                        {item}
                      </Anchor>
                    )
                  })}
                </Breadcrumbs>
                {currentPath && (
                  <Button
                    variant="subtle"
                    leftSection={<ArrowLeft size={16} />}
                    onClick={navigateBack}
                    size={isMobile ? 'xs' : 'sm'}
                  >
                    Back
                  </Button>
                )}
              </Group>
            </Paper>

            {error && (
              <Notification color="red" title="Error" onClose={() => setError(null)}>
                {error}
              </Notification>
            )}

            <Paper withBorder p="md" radius="md">
              <Stack>
                <Dropzone
                  onDrop={onDrop}
                  loading={uploading}
                  maxSize={config.maxSizeBytes}
                  disabled={config.readOnly}
                  styles={{
                    root: {
                      border: `${rem(1)} dashed var(--mantine-color-default-border)`,
                      borderRadius: 'var(--mantine-radius-md)',
                      padding: isMobile ? rem(8) : rem(16),
                      backgroundColor: 'var(--mantine-color-default-hover)',
                      cursor: 'pointer',
                      '&:hover': {
                        backgroundColor: 'var(--mantine-color-default)',
                      },
                    },
                  }}
                >
                  <Group justify="center" gap={isMobile ? 'xs' : 'lg'} mih={isMobile ? 60 : 80} style={{ pointerEvents: 'none' }}>
                    <Dropzone.Accept>
                      <Upload size={isMobile ? 24 : 32} color="var(--mantine-color-violet-6)" />
                    </Dropzone.Accept>
                    <Dropzone.Reject>
                      <File size={isMobile ? 24 : 32} color="var(--mantine-color-red-6)" />
                    </Dropzone.Reject>
                    <Dropzone.Idle>
                      <Upload size={isMobile ? 24 : 32} color="var(--mantine-color-dimmed)" />
                    </Dropzone.Idle>

                    <div>
                      <Text size={isMobile ? 'sm' : 'md'} fw={600} ta={isMobile ? 'left' : 'center'}>
                        {isMobile ? 'Tap to upload files' : 'Drag files here or click to select files'}
                      </Text>
                      {!isMobile && (
                        <Text size="xs" c="dimmed" ta="center">
                          Max {formatBytes(config.maxSizeBytes)} per request
                        </Text>
                      )}
                    </div>
                  </Group>
                </Dropzone>
                {uploading && (
                  <Stack gap={4}>
                    <Group justify="space-between">
                      <Text size="xs" fw={500}>Uploading...</Text>
                      <Text size="xs" fw={500}>{uploadProgress}%</Text>
                    </Group>
                    <Progress value={uploadProgress} size="sm" color="violet" animated />
                  </Stack>
                )}
              </Stack>
            </Paper>

            {isMobile ? (
              <Stack gap="sm">{mobileCards}</Stack>
            ) : (
              <Table verticalSpacing="sm" highlightOnHover>
                <Table.Thead>
                  <Table.Tr>
                    <Table.Th>Name</Table.Th>
                    <Table.Th>Size</Table.Th>
                    <Table.Th>Modified</Table.Th>
                    <Table.Th w={80}>Action</Table.Th>
                  </Table.Tr>
                </Table.Thead>
                <Table.Tbody>{rows}</Table.Tbody>
              </Table>
            )}
          </Stack>
        </Container>
      </AppShell.Main>
    </AppShell>
  )
}

export default App
