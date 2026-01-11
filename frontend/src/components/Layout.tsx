import { Outlet, Link, useLocation } from 'react-router-dom'
import {
  Navbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
} from '@heroui/react'
import { ThemeToggle } from './ThemeToggle'

export default function Layout() {
  const location = useLocation()

  const isActive = (path: string) => {
    if (path === '/') {
      return location.pathname === '/'
    }
    return location.pathname.startsWith(path)
  }

  return (
    <div className="min-h-screen bg-background">
      <Navbar maxWidth="xl" isBordered>
        <NavbarBrand>
          <Link to="/" className="font-bold text-inherit text-xl">
            Dict Hub
          </Link>
        </NavbarBrand>

        <NavbarContent className="gap-4" justify="center">
          <NavbarItem isActive={isActive('/')}>
            <Link
              to="/"
              className={
                isActive('/')
                  ? 'text-primary font-medium'
                  : 'text-foreground hover:text-primary'
              }
            >
              首页
            </Link>
          </NavbarItem>
          <NavbarItem isActive={isActive('/settings')}>
            <Link
              to="/settings"
              className={
                isActive('/settings')
                  ? 'text-primary font-medium'
                  : 'text-foreground hover:text-primary'
              }
            >
              设置
            </Link>
          </NavbarItem>
          <NavbarItem isActive={isActive('/history')}>
            <Link
              to="/history"
              className={
                isActive('/history')
                  ? 'text-primary font-medium'
                  : 'text-foreground hover:text-primary'
              }
            >
              历史
            </Link>
          </NavbarItem>
        </NavbarContent>

        <NavbarContent justify="end">
          <NavbarItem>
            <ThemeToggle />
          </NavbarItem>
        </NavbarContent>
      </Navbar>

      <main className="container mx-auto max-w-6xl p-6">
        <Outlet />
      </main>
    </div>
  )
}
