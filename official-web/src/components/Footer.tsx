import { motion } from 'framer-motion';

const footerLinks = [
  { label: 'GitHub', href: 'https://github.com/Davied-H/dict-hub' },
  { label: '问题反馈', href: 'https://github.com/Davied-H/dict-hub/issues' },
  { label: '许可证', href: 'https://github.com/Davied-H/dict-hub/blob/main/LICENSE' },
];

export default function Footer() {
  return (
    <footer className="relative py-12 bg-slate-50 border-t border-slate-200">
      <div className="max-w-7xl mx-auto px-6">
        <div className="flex flex-col md:flex-row justify-between items-center gap-6">
          <motion.div
            className="flex items-center gap-3"
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
          >
            <div className="w-8 h-8 bg-gradient-primary rounded-lg flex items-center justify-center text-white">
              📚
            </div>
            <span className="text-slate-600 text-sm">
              © 2026 Dict-Hub. Made with ❤️
            </span>
          </motion.div>

          <motion.div
            className="flex gap-8"
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.1 }}
          >
            {footerLinks.map((link) => (
              <motion.a
                key={link.label}
                href={link.href}
                target="_blank"
                rel="noopener noreferrer"
                className="text-slate-600 text-sm hover:text-primary transition-colors relative group"
                whileHover={{ y: -2 }}
              >
                {link.label}
                <span className="absolute -bottom-1 left-0 w-0 h-0.5 bg-primary group-hover:w-full transition-all duration-300" />
              </motion.a>
            ))}
          </motion.div>
        </div>
      </div>
    </footer>
  );
}
