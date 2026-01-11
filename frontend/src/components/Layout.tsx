import { useState } from 'react'
import { Outlet, Link, useLocation } from 'react-router-dom'
import {
  Navbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
  NavbarMenuToggle,
  NavbarMenu,
  NavbarMenuItem,
} from '@heroui/react'
import { ThemeToggle } from './ThemeToggle'

const menuItems = [
  { name: '首页', path: '/' },
  { name: '设置', path: '/settings' },
  { name: '历史', path: '/history' },
]

export default function Layout() {
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const location = useLocation()

  const isActive = (path: string) => {
    if (path === '/') {
      return location.pathname === '/'
    }
    return location.pathname.startsWith(path)
  }

  return (
    <div className="min-h-screen bg-background">
      <Navbar maxWidth="xl" isBordered isMenuOpen={isMenuOpen} onMenuOpenChange={setIsMenuOpen}>
        {/* 移动端汉堡菜单按钮 */}
        <NavbarContent className="sm:hidden" justify="start">
          <NavbarMenuToggle aria-label={isMenuOpen ? '关闭菜单' : '打开菜单'} />
        </NavbarContent>

        {/* Logo - 移动端居中，桌面端靠左 */}
        <NavbarBrand className="sm:flex-grow-0">
          <Link to="/" className="font-bold text-inherit text-xl">
            Dict Hub
          </Link>
        </NavbarBrand>

        {/* 桌面端导航链接 */}
        <NavbarContent className="hidden sm:flex gap-4" justify="center">
          {menuItems.map((item) => (
            <NavbarItem key={item.path} isActive={isActive(item.path)}>
              <Link
                to={item.path}
                className={
                  isActive(item.path)
                    ? 'text-primary font-medium'
                    : 'text-foreground hover:text-primary'
                }
              >
                {item.name}
              </Link>
            </NavbarItem>
          ))}
        </NavbarContent>

        {/* 主题切换按钮 */}
        <NavbarContent justify="end">
          <NavbarItem>
            <ThemeToggle />
          </NavbarItem>
        </NavbarContent>

        {/* 移动端下拉菜单 */}
        <NavbarMenu>
          {menuItems.map((item) => (
            <NavbarMenuItem key={item.path}>
              <Link
                to={item.path}
                className={`w-full py-2 block ${
                  isActive(item.path)
                    ? 'text-primary font-medium'
                    : 'text-foreground hover:text-primary'
                }`}
                onClick={() => setIsMenuOpen(false)}
              >
                {item.name}
              </Link>
            </NavbarMenuItem>
          ))}
        </NavbarMenu>
      </Navbar>

      <main className="container mx-auto max-w-6xl px-4 py-4 sm:p-6">
        <Outlet />
      </main>
    </div>
  )
}
