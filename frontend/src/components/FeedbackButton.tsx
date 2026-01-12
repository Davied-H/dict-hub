import {
  Button,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  Card,
  CardBody,
  useDisclosure,
} from '@heroui/react'

const GITHUB_REPO = 'https://github.com/Davied-H/dict-hub'

interface FeedbackType {
  id: string
  title: string
  description: string
  icon: string
  label: string
  template: string
}

const feedbackTypes: FeedbackType[] = [
  {
    id: 'bug',
    title: '问题反馈',
    description: '报告 Bug 或使用问题',
    icon: '🐛',
    label: 'bug',
    template: `**问题描述**
<!-- 请描述你遇到的问题 -->

**复现步骤**
1.
2.
3.

**预期行为**
<!-- 你期望发生什么？ -->

**实际行为**
<!-- 实际发生了什么？ -->

**环境信息**
- 浏览器:
- 操作系统: `,
  },
  {
    id: 'feature',
    title: '功能建议',
    description: '提出新功能或改进建议',
    icon: '💡',
    label: 'enhancement',
    template: `**功能描述**
<!-- 请描述你希望添加的功能 -->

**使用场景**
<!-- 这个功能会在什么场景下使用？ -->

**可能的实现方式**
<!-- 如果有想法，请描述可能的实现方式 -->`,
  },
]

function openGitHubIssue(feedbackType: FeedbackType) {
  const params = new URLSearchParams({
    labels: feedbackType.label,
    body: feedbackType.template,
  })
  const url = `${GITHUB_REPO}/issues/new?${params.toString()}`
  window.open(url, '_blank')
}

interface FeedbackButtonProps {
  variant?: 'navbar' | 'menu'
  onClose?: () => void
}

export function FeedbackButton({ variant = 'navbar', onClose }: FeedbackButtonProps) {
  const { isOpen, onOpen, onOpenChange } = useDisclosure()

  const handleSelectType = (type: FeedbackType) => {
    openGitHubIssue(type)
    onOpenChange()
    onClose?.()
  }

  if (variant === 'menu') {
    return (
      <>
        <button
          onClick={onOpen}
          className="w-full py-2 block text-left text-foreground hover:text-primary"
        >
          反馈
        </button>
        <FeedbackModal
          isOpen={isOpen}
          onOpenChange={onOpenChange}
          onSelectType={handleSelectType}
        />
      </>
    )
  }

  return (
    <>
      <Button
        variant="light"
        size="sm"
        onPress={onOpen}
        className="text-foreground"
      >
        反馈
      </Button>
      <FeedbackModal
        isOpen={isOpen}
        onOpenChange={onOpenChange}
        onSelectType={handleSelectType}
      />
    </>
  )
}

interface FeedbackModalProps {
  isOpen: boolean
  onOpenChange: () => void
  onSelectType: (type: FeedbackType) => void
}

function FeedbackModal({ isOpen, onOpenChange, onSelectType }: FeedbackModalProps) {
  return (
    <Modal isOpen={isOpen} onOpenChange={onOpenChange} placement="center">
      <ModalContent>
        <ModalHeader className="flex flex-col gap-1">
          选择反馈类型
        </ModalHeader>
        <ModalBody className="pb-6">
          <div className="grid gap-3">
            {feedbackTypes.map((type) => (
              <Card
                key={type.id}
                isPressable
                onPress={() => onSelectType(type)}
                className="hover:bg-default-100"
              >
                <CardBody className="flex flex-row items-center gap-4">
                  <span className="text-2xl">{type.icon}</span>
                  <div>
                    <p className="font-medium">{type.title}</p>
                    <p className="text-sm text-default-500">{type.description}</p>
                  </div>
                </CardBody>
              </Card>
            ))}
          </div>
        </ModalBody>
      </ModalContent>
    </Modal>
  )
}
