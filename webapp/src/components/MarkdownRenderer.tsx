'use client';

import Markdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { Noto_Sans_Mono } from 'next/font/google';
import classNames from 'classnames';

const mono = Noto_Sans_Mono({
  subsets: ['latin'],
  display: 'swap',
});

export default function MarkdownRenderer({
  className,
  loading,
  content,
}: {
  className?: string;
  loading?: boolean;
  content: string;
}) {
  return (
    <Markdown
      className={classNames(className, {
        'animate-pulse': loading,
      })}
      remarkPlugins={[remarkGfm]}
      components={{
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        code: ({ children, className, ref, ...restProps }) => {
          const match = /language-(\w+)/.exec(className || '');
          return match ? (
            <SyntaxHighlighter
              {...restProps}
              language={match[1]}
              style={vscDarkPlus}
              customStyle={{
                padding: '0.5rem',
                margin: '0rem',
                overflow: 'auto',
                backgroundColor: 'transparent',
                ...mono.style,
                fontWeight: 500,
              }}
            >
              {String(children).replace(/\n$/, '')}
            </SyntaxHighlighter>
          ) : (
            <span
              className="m-0 whitespace-break-spaces rounded-md bg-[#656c7633] px-1.5 py-0.5"
              style={{
                fontSize: '85%',
                ...mono.style,
                fontWeight: 500,
              }}
            >
              {children}
            </span>
          );
        },
      }}
    >
      {content}
    </Markdown>
  );
}
