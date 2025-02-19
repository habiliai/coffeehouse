'use client';

import {
  Carousel,
  CarouselApi,
  CarouselContent,
  CarouselItem,
} from '@/components/ui/carousel';
import { useEffect, useState } from 'react';
import classNames from 'classnames';

export type CarouselItem = {
  value: string;
  name: string;
};

export function VerticalCarousel({
  className,
  itemClassName,
  items,
  selectedValue,
  onClick,
}: {
  className?: string;
  itemClassName?: string;
  items: CarouselItem[];
  selectedValue?: string;
  onClick?: (value: string) => void;
}) {
  const [carouselApi, setCarouselApi] = useState<CarouselApi>();

  useEffect(() => {
    if (!carouselApi || selectedValue) return;

    const interval = setInterval(() => {
      carouselApi.scrollNext();
    }, 3000);
    return () => clearInterval(interval);
  }, [carouselApi, selectedValue]);

  useEffect(() => {
    if (!carouselApi || !selectedValue) return;

    const index = items.findIndex((item) => item.value === selectedValue);
    if (index !== -1) {
      carouselApi.scrollTo(index - 1);
    }
  }, [carouselApi, selectedValue, items]);

  return (
    <div className={className}>
      <Carousel
        setApi={setCarouselApi}
        opts={{
          align: 'start',
          loop: true,
        }}
        orientation="vertical"
        className="w-full"
      >
        <CarouselContent className="-mt-1 h-[3.685rem] lg:h-[12.5rem]">
          {items.map((item, index) => (
            <CarouselItem
              key={`carousel-item-${index}`}
              className="basis-1/3 cursor-pointer pt-1"
              onClick={() => onClick?.(item.value)}
            >
              <div
                className={classNames(
                  'flex h-full items-center p-1 text-3xl font-semibold',
                  itemClassName,
                  {
                    'bg-gray-200': selectedValue === item.value,
                    'opacity-50': selectedValue !== item.value,
                  },
                )}
              >
                <span className="line-clamp-1">{item.name}</span>
              </div>
            </CarouselItem>
          ))}
        </CarouselContent>
      </Carousel>
    </div>
  );
}
