import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Card,
  CardBody,
  Table,
  TableHeader,
  TableColumn,
  TableBody,
  TableRow,
  TableCell,
  Button,
  Pagination,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Spinner,
  Chip,
} from '@heroui/react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { historyApi } from '@/api'
import { useSearchStore } from '@/stores'

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export default function History() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { setKeyword } = useSearchStore()
  const [page, setPage] = useState(1)
  const [isClearModalOpen, setIsClearModalOpen] = useState(false)
  const pageSize = 20

  const { data, isLoading, error } = useQuery({
    queryKey: ['history', page],
    queryFn: async () => {
      const res = await historyApi.list(page, pageSize)
      return res.data
    },
  })

  const clearMutation = useMutation({
    mutationFn: () => historyApi.clear(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['history'] })
      setIsClearModalOpen(false)
      setPage(1)
    },
  })

  const handleExport = () => {
    const link = document.createElement('a')
    link.href = historyApi.exportUrl
    link.download = `history_${Date.now()}.csv`
    link.click()
  }

  const handleRowClick = (word: string) => {
    setKeyword(word)
    navigate('/')
  }

  const totalPages = data?.meta
    ? Math.ceil(data.meta.total / data.meta.pageSize)
    : 1
  const historyList = data?.data || []

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-16">
        <Spinner size="lg" />
      </div>
    )
  }

  if (error) {
    return (
      <Card>
        <CardBody className="text-danger">加载历史记录失败</CardBody>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">搜索历史</h1>
        <div className="flex gap-2">
          <Button
            variant="flat"
            onPress={handleExport}
            isDisabled={historyList.length === 0}
          >
            导出 CSV
          </Button>
          <Button
            color="danger"
            variant="flat"
            onPress={() => setIsClearModalOpen(true)}
            isDisabled={historyList.length === 0}
          >
            清空历史
          </Button>
        </div>
      </div>

      <Card>
        <CardBody>
          {historyList.length === 0 ? (
            <div className="text-center py-8 text-default-500">
              暂无搜索历史
            </div>
          ) : (
            <>
              <Table
                aria-label="搜索历史"
                selectionMode="single"
                onRowAction={(key) => {
                  const item = historyList.find((h) => h.id === Number(key))
                  if (item) handleRowClick(item.word)
                }}
              >
                <TableHeader>
                  <TableColumn>单词</TableColumn>
                  <TableColumn>结果</TableColumn>
                  <TableColumn>响应时间</TableColumn>
                  <TableColumn>时间</TableColumn>
                </TableHeader>
                <TableBody>
                  {historyList.map((item) => (
                    <TableRow key={item.id} className="cursor-pointer">
                      <TableCell>
                        <span className="font-medium">{item.word}</span>
                      </TableCell>
                      <TableCell>
                        <Chip
                          size="sm"
                          color={item.found ? 'success' : 'default'}
                          variant="flat"
                        >
                          {item.found ? '找到' : '未找到'}
                        </Chip>
                      </TableCell>
                      <TableCell>{item.responseTime}ms</TableCell>
                      <TableCell>{formatDate(item.createdAt)}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>

              {totalPages > 1 && (
                <div className="flex justify-center mt-4">
                  <Pagination
                    total={totalPages}
                    page={page}
                    onChange={setPage}
                    showControls
                  />
                </div>
              )}
            </>
          )}
        </CardBody>
      </Card>

      <Modal
        isOpen={isClearModalOpen}
        onClose={() => setIsClearModalOpen(false)}
      >
        <ModalContent>
          <ModalHeader>确认清空</ModalHeader>
          <ModalBody>确定要清空所有搜索历史吗？此操作无法撤销。</ModalBody>
          <ModalFooter>
            <Button variant="light" onPress={() => setIsClearModalOpen(false)}>
              取消
            </Button>
            <Button
              color="danger"
              onPress={() => clearMutation.mutate()}
              isLoading={clearMutation.isPending}
            >
              清空
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </div>
  )
}
