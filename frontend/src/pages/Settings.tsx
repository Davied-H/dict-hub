import { useState } from 'react'
import {
  Card,
  CardHeader,
  CardBody,
  Table,
  TableHeader,
  TableColumn,
  TableBody,
  TableRow,
  TableCell,
  Switch,
  Button,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Input,
  Tabs,
  Tab,
  Spinner,
} from '@heroui/react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { dictionariesApi } from '@/api'
import { useThemeStore, type Theme } from '@/stores'
import type { DictSource } from '@/types'

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

function DictManager() {
  const queryClient = useQueryClient()
  const [isAddModalOpen, setIsAddModalOpen] = useState(false)
  const [newDictPath, setNewDictPath] = useState('')
  const [deleteConfirmId, setDeleteConfirmId] = useState<number | null>(null)

  const { data, isLoading, error } = useQuery({
    queryKey: ['dictionaries'],
    queryFn: async () => {
      const res = await dictionariesApi.list()
      return res.data.data || []
    },
  })

  const toggleMutation = useMutation({
    mutationFn: (id: number) => dictionariesApi.toggle(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['dictionaries'] })
    },
  })

  const addMutation = useMutation({
    mutationFn: (path: string) => dictionariesApi.add(path),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['dictionaries'] })
      setIsAddModalOpen(false)
      setNewDictPath('')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => dictionariesApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['dictionaries'] })
      setDeleteConfirmId(null)
    },
  })

  const reorderMutation = useMutation({
    mutationFn: (orders: { id: number; sort_order: number }[]) =>
      dictionariesApi.reorder(orders),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['dictionaries'] })
    },
  })

  const handleMove = (dict: DictSource, direction: 'up' | 'down') => {
    if (!data) return
    const sorted = [...data].sort((a, b) => a.sort_order - b.sort_order)
    const currentIndex = sorted.findIndex((d) => d.id === dict.id)
    const targetIndex = direction === 'up' ? currentIndex - 1 : currentIndex + 1

    if (targetIndex < 0 || targetIndex >= sorted.length) return

    const newOrders = sorted.map((d, i) => {
      if (i === currentIndex) return { id: d.id, sort_order: targetIndex }
      if (i === targetIndex) return { id: d.id, sort_order: currentIndex }
      return { id: d.id, sort_order: i }
    })

    reorderMutation.mutate(newOrders)
  }

  const sortedData = data
    ? [...data].sort((a, b) => a.sort_order - b.sort_order)
    : []

  if (isLoading) {
    return (
      <Card>
        <CardBody className="flex items-center justify-center py-8">
          <Spinner size="lg" />
        </CardBody>
      </Card>
    )
  }

  if (error) {
    return (
      <Card>
        <CardBody className="text-danger">加载字典列表失败</CardBody>
      </Card>
    )
  }

  return (
    <>
      <Card>
        <CardHeader className="flex justify-between items-center">
          <h2 className="text-xl font-semibold">字典管理</h2>
          <Button color="primary" onPress={() => setIsAddModalOpen(true)}>
            添加字典
          </Button>
        </CardHeader>
        <CardBody>
          {sortedData.length === 0 ? (
            <div className="text-center py-8 text-default-500">
              暂无字典，请添加
            </div>
          ) : (
            <Table aria-label="字典列表">
              <TableHeader>
                <TableColumn>名称</TableColumn>
                <TableColumn>词条数</TableColumn>
                <TableColumn>大小</TableColumn>
                <TableColumn>状态</TableColumn>
                <TableColumn>排序</TableColumn>
                <TableColumn>操作</TableColumn>
              </TableHeader>
              <TableBody>
                {sortedData.map((dict, index) => (
                  <TableRow key={dict.id}>
                    <TableCell>
                      <div>
                        <div className="font-medium">
                          {dict.title || dict.name}
                        </div>
                        {dict.description && (
                          <div className="text-xs text-default-400 truncate max-w-xs">
                            {dict.description}
                          </div>
                        )}
                      </div>
                    </TableCell>
                    <TableCell>{dict.word_count.toLocaleString()}</TableCell>
                    <TableCell>{formatFileSize(dict.file_size)}</TableCell>
                    <TableCell>
                      <Switch
                        isSelected={dict.enabled}
                        onValueChange={() => toggleMutation.mutate(dict.id)}
                        isDisabled={toggleMutation.isPending}
                        size="sm"
                      />
                    </TableCell>
                    <TableCell>
                      <div className="flex gap-1">
                        <Button
                          isIconOnly
                          size="sm"
                          variant="light"
                          isDisabled={index === 0 || reorderMutation.isPending}
                          onPress={() => handleMove(dict, 'up')}
                        >
                          ↑
                        </Button>
                        <Button
                          isIconOnly
                          size="sm"
                          variant="light"
                          isDisabled={
                            index === sortedData.length - 1 ||
                            reorderMutation.isPending
                          }
                          onPress={() => handleMove(dict, 'down')}
                        >
                          ↓
                        </Button>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Button
                        size="sm"
                        color="danger"
                        variant="light"
                        onPress={() => setDeleteConfirmId(dict.id)}
                      >
                        删除
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardBody>
      </Card>

      <Modal isOpen={isAddModalOpen} onClose={() => setIsAddModalOpen(false)}>
        <ModalContent>
          <ModalHeader>添加字典</ModalHeader>
          <ModalBody>
            <Input
              label="字典路径"
              placeholder="输入 MDX 文件的完整路径"
              value={newDictPath}
              onChange={(e) => setNewDictPath(e.target.value)}
            />
          </ModalBody>
          <ModalFooter>
            <Button variant="light" onPress={() => setIsAddModalOpen(false)}>
              取消
            </Button>
            <Button
              color="primary"
              onPress={() => addMutation.mutate(newDictPath)}
              isLoading={addMutation.isPending}
              isDisabled={!newDictPath.trim()}
            >
              添加
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>

      <Modal
        isOpen={deleteConfirmId !== null}
        onClose={() => setDeleteConfirmId(null)}
      >
        <ModalContent>
          <ModalHeader>确认删除</ModalHeader>
          <ModalBody>确定要删除这个字典吗？此操作无法撤销。</ModalBody>
          <ModalFooter>
            <Button variant="light" onPress={() => setDeleteConfirmId(null)}>
              取消
            </Button>
            <Button
              color="danger"
              onPress={() =>
                deleteConfirmId && deleteMutation.mutate(deleteConfirmId)
              }
              isLoading={deleteMutation.isPending}
            >
              删除
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  )
}

function ThemeSettings() {
  const { theme, setTheme } = useThemeStore()

  return (
    <Card>
      <CardHeader>
        <h2 className="text-xl font-semibold">主题设置</h2>
      </CardHeader>
      <CardBody>
        <Tabs
          selectedKey={theme}
          onSelectionChange={(key) => setTheme(key as Theme)}
          aria-label="主题选择"
        >
          <Tab key="light" title="亮色" />
          <Tab key="dark" title="暗色" />
          <Tab key="system" title="跟随系统" />
        </Tabs>
      </CardBody>
    </Card>
  )
}

export default function Settings() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">设置</h1>
      <DictManager />
      <ThemeSettings />
    </div>
  )
}
